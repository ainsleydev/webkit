package appdef

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/afero"
)

var (
	// validate is the singleton validator instance with custom validators registered.
	validate *validator.Validate
)

func init() {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("lowercase", validateLowercase)
	_ = validate.RegisterValidation("alphanumdash", validateAlphanumDash)
}

// Validate performs comprehensive validation on the Definition,
// collecting all errors before returning. Returns nil if validation passes.
func (d *Definition) Validate(fs afero.Fs) []error {
	var errs []error

	// Struct validation with go-playground/validator
	errs = append(errs, validateStruct(d)...)

	// Business logic validation
	errs = append(errs, d.validateDomains()...)
	errs = append(errs, d.validateAppPaths(fs)...)
	errs = append(errs, d.validateTerraformManagedVMs()...)
	errs = append(errs, d.validateEnvReferences()...)

	// Return nil if no errors
	if len(errs) == 0 {
		return nil
	}

	return errs
}

// validateStruct validates a Definition struct using go-playground/validator.
func validateStruct(def *Definition) []error {
	err := validate.Struct(def)
	if err == nil {
		return nil
	}

	// Convert validator errors to []error
	var errs []error
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, e := range validationErrs {
			errs = append(errs, fmt.Errorf("%s: validation failed on '%s' tag", e.Field(), e.Tag()))
		}
	}

	return errs
}

// validateLowercase checks if a string contains only lowercase characters.
func validateLowercase(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, r := range value {
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// validateAlphanumDash validates that a string matches the pattern ^[a-z][a-z0-9-]*$
// (starts with a lowercase letter, followed by lowercase letters, numbers, or hyphens).
func validateAlphanumDash(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9-]*$`, value)
	return matched
}

// validateDomains ensures that domain names do not contain protocol prefixes.
func (d *Definition) validateDomains() []error {
	var errs []error

	for _, app := range d.Apps {
		for _, domain := range app.Domains {
			if strings.Contains(domain.Name, "://") {
				errs = append(errs, fmt.Errorf(
					"app %q: domain %q should not contain protocol prefix (e.g., 'https://')",
					app.Name,
					domain.Name,
				))
			}
		}
	}

	return errs
}

// validateAppPaths ensures that all app paths exist on the filesystem.
func (d *Definition) validateAppPaths(fs afero.Fs) []error {
	var errs []error

	for _, app := range d.Apps {
		if app.Path == "" {
			continue
		}

		exists, err := afero.DirExists(fs, app.Path)
		if err != nil {
			errs = append(errs, fmt.Errorf(
				"app %q: error checking path %q: %w",
				app.Name,
				app.Path,
				err,
			))
			continue
		}

		if !exists {
			errs = append(errs, fmt.Errorf(
				"app %q: path %q does not exist",
				app.Name,
				app.Path,
			))
		}
	}

	return errs
}

// validateTerraformManagedVMs ensures that terraform-managed VM apps
// have at least one domain configured.
func (d *Definition) validateTerraformManagedVMs() []error {
	var errs []error

	for _, app := range d.Apps {
		// Check if this is a terraform-managed VM/app
		if app.IsTerraformManaged() && (app.Infra.Type == "vm" || app.Infra.Type == "app") {
			if len(app.Domains) == 0 {
				errs = append(errs, fmt.Errorf(
					"app %q: terraform-managed VM/app must have at least one domain configured",
					app.Name,
				))
			}
		}
	}

	return errs
}

// validateEnvReferences ensures that all environment variable resource
// references point to valid resources and outputs.
func (d *Definition) validateEnvReferences() []error {
	var errs []error

	// Build a map of resource names to their types for quick lookup
	resourceMap := make(map[string]ResourceType)
	for _, res := range d.Resources {
		resourceMap[res.Name] = res.Type
	}

	// Validate shared env references
	errs = append(errs, d.validateEnvVarReferences("shared", d.Shared.Env, resourceMap)...)

	// Validate each app's env references
	for _, app := range d.Apps {
		errs = append(errs, d.validateEnvVarReferences(
			fmt.Sprintf("app %q", app.Name),
			app.Env,
			resourceMap,
		)...)
	}

	return errs
}

// validateEnvVarReferences validates environment variable references for a
// given context (shared or app-specific).
func (d *Definition) validateEnvVarReferences(
	context string,
	env Environment,
	resourceMap map[string]ResourceType,
) []error {
	var errs []error

	// Walk through all env vars
	err := env.WalkE(func(entry EnvWalkEntry) error {
		// Only validate resource references
		if entry.Source != EnvSourceResource {
			return nil
		}

		// Parse the resource reference
		resourceName, outputName, ok := ParseResourceReference(entry.Value)
		if !ok {
			errs = append(errs, fmt.Errorf(
				"%s: env var %q in %s has invalid resource reference format %q (expected 'resource_name.output_name')",
				context,
				entry.Key,
				entry.Environment,
				entry.Value,
			))
			return nil
		}

		// Check if resource exists
		resourceType, exists := resourceMap[resourceName]
		if !exists {
			errs = append(errs, fmt.Errorf(
				"%s: env var %q in %s references non-existent resource %q",
				context,
				entry.Key,
				entry.Environment,
				resourceName,
			))
			return nil
		}

		// Check if output is valid for this resource type
		validOutputs := resourceType.Outputs()
		if validOutputs != nil {
			outputValid := false
			for _, validOutput := range validOutputs {
				if outputName == validOutput {
					outputValid = true
					break
				}
			}

			if !outputValid {
				errs = append(errs, fmt.Errorf(
					"%s: env var %q in %s references invalid output %q for resource %q (type: %s). Valid outputs: %v",
					context,
					entry.Key,
					entry.Environment,
					outputName,
					resourceName,
					resourceType,
					validOutputs,
				))
			}
		}

		return nil
	})
	if err != nil {
		errs = append(errs, fmt.Errorf("walking env variables: %w", err))
	}

	return errs
}
