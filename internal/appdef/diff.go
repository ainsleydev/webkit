package appdef

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/appdef/types"
	"github.com/ainsleydev/webkit/internal/util/executil"
)

type (
	// ChangeAnalysis represents the result of comparing app.json versions.
	ChangeAnalysis struct {
		// Skip indicates whether Terraform apply should be skipped.
		Skip bool `json:"skip"`

		// Reason explains why the decision was made.
		Reason string `json:"reason"`

		// ChangedApps lists apps that have changes.
		ChangedApps []AppChange `json:"changed_apps,omitempty"`
	}

	// AppChange represents changes detected for a specific app.
	AppChange struct {
		// Name of the app.
		Name string `json:"name"`

		// InfraChanged indicates non-env infrastructure changes.
		InfraChanged bool `json:"infra_changed"`

		// EnvChanged indicates environment variable changes.
		EnvChanged bool `json:"env_changed"`

		// PlatformType is the type of platform (vm, container, serverless).
		PlatformType string `json:"platform_type"`

		// PlatformProvider is the provider (digitalocean, aws, etc).
		PlatformProvider string `json:"platform_provider"`
	}
)

// Compare compares two app.json definitions and determines if Terraform apply is needed.
//
// It performs a hierarchical comparison:
// 1. If definitions are identical → skip Terraform
// 2. If infrastructure config changed (non-env fields) → run Terraform
// 3. If only env vars changed → analyse per app:
//   - VM or non-DigitalOcean apps → run Terraform
//   - DigitalOcean container apps with env value changes → run Terraform
//   - DigitalOcean container apps with no actual env changes → skip (drift only)
func Compare(current, previous *Definition) ChangeAnalysis {
	// Ignore unexported fields in OrderedMap.
	ignoreOrderedMap := cmpopts.IgnoreUnexported(types.OrderedMap[Command, CommandSpec]{})

	// Quick check: if definitions are identical, skip.
	if cmp.Equal(current, previous, ignoreOrderedMap) {
		return ChangeAnalysis{
			Skip:   true,
			Reason: "app.json unchanged",
		}
	}

	// Check non-env fields (infrastructure config).
	// Ignore env fields in both shared and apps.
	ignoreEnvs := cmpopts.IgnoreFields(Environment{}, "Dev", "Staging", "Production")
	ignoreSharedEnv := cmpopts.IgnoreFields(Shared{}, "Env")
	ignoreAppEnv := cmpopts.IgnoreFields(App{}, "Env")

	if !cmp.Equal(current, previous, ignoreOrderedMap, ignoreEnvs, ignoreSharedEnv, ignoreAppEnv) {
		return ChangeAnalysis{
			Skip:   false,
			Reason: "Infrastructure config changed (domains/sizes/regions/resources/etc)",
		}
	}

	// Only env vars changed - analyse per app.
	changedApps := analyseEnvChanges(current, previous)

	// Check if any changes require terraform.
	for _, app := range changedApps {
		// VM or non-DO container apps always need terraform.
		if app.PlatformType != "container" || app.PlatformProvider != "digitalocean" {
			return ChangeAnalysis{
				Skip:        false,
				Reason:      "VM or non-DigitalOcean container app env changes detected",
				ChangedApps: changedApps,
			}
		}

		// If DO container env values actually changed, need terraform.
		if app.EnvChanged {
			return ChangeAnalysis{
				Skip:        false,
				Reason:      "DigitalOcean container app env values changed",
				ChangedApps: changedApps,
			}
		}
	}

	// Only DO container env vars touched, but values identical (drift only).
	return ChangeAnalysis{
		Skip:        true,
		Reason:      "Only DigitalOcean container env vars touched, but values unchanged (drift only)",
		ChangedApps: changedApps,
	}
}

// analyseEnvChanges compares environment variables for each app.
func analyseEnvChanges(current, previous *Definition) []AppChange {
	changes := make([]AppChange, 0)

	// Create a map of previous apps for easier lookup.
	previousApps := make(map[string]*App)
	for i := range previous.Apps {
		previousApps[previous.Apps[i].Name] = &previous.Apps[i]
	}

	// Compare each current app with its previous version.
	for i := range current.Apps {
		currentApp := &current.Apps[i]
		previousApp, exists := previousApps[currentApp.Name]

		// New app added - infrastructure change.
		if !exists {
			continue
		}

		// Check if env vars changed using cmp.
		envChanged := !cmp.Equal(currentApp.Env, previousApp.Env)

		if envChanged {
			changes = append(changes, AppChange{
				Name:             currentApp.Name,
				InfraChanged:     false,
				EnvChanged:       envChanged,
				PlatformType:     currentApp.Infra.Type,
				PlatformProvider: currentApp.Infra.Provider.String(),
			})
		}
	}

	return changes
}

// LoadFromGit loads an app.json definition from a specific git reference.
//
// The ref parameter can be any valid git reference (e.g., "HEAD~1", "origin/main", commit SHA).
//
// Returns an error if git operations fail or if the file cannot be parsed.
func LoadFromGit(ctx context.Context, ref string) (*Definition, error) {
	cmd := executil.NewCommand("git", "show", ref+":"+JsonFileName)
	result, err := executil.DefaultRunner().Run(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "executing git show")
	}

	def, err := Parse([]byte(result.Output))
	if err != nil {
		return nil, errors.Wrap(err, "parsing app.json from git")
	}

	return def, nil
}
