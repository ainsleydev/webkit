package appdef

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Validate performs comprehensive validation on the Definition,
// collecting all errors before returning. Returns nil if validation passes.
func (d *Definition) Validate(fs afero.Fs) []error {
	var errs []error

	// Validate domains
	errs = append(errs, d.validateDomains()...)

	// Validate app paths
	errs = append(errs, d.validateAppPaths(fs)...)

	// Validate terraform-managed VMs have domains
	errs = append(errs, d.validateTerraformManagedVMs()...)

	// Validate env references
	errs = append(errs, d.validateEnvReferences()...)

	// Return nil if no errors
	if len(errs) == 0 {
		return nil
	}

	return errs
}

// validateDomains ensures that domain names do not contain protocol prefixes.
func (d *Definition) validateDomains() []error {
	var errs []error

	for _, app := range d.Apps {
		for _, domain := range app.Domains {
			if strings.Contains(domain.Name, "://") {
				errs = append(errs, errors.Errorf(
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
			errs = append(errs, errors.Wrapf(
				err,
				"app %q: error checking path %q",
				app.Name,
				app.Path,
			))
			continue
		}

		if !exists {
			errs = append(errs, errors.Errorf(
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
				errs = append(errs, errors.Errorf(
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
			errs = append(errs, errors.Errorf(
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
			errs = append(errs, errors.Errorf(
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
				errs = append(errs, errors.Errorf(
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
		errs = append(errs, errors.Wrap(err, "walking env variables"))
	}

	return errs
}
