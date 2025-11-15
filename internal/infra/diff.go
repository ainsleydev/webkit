package infra

import (
	"context"
	"os"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/appdef"
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

// Diff compares current app.json with a previous version and determines
// if Terraform apply is needed.
//
// The baseRef parameter specifies the git reference to compare against
// (e.g., "HEAD~1", "origin/main").
//
// Returns an error if git operations fail or if app.json cannot be parsed.
func Diff(ctx context.Context, baseRef string) (ChangeAnalysis, error) {
	// Load current app.json.
	currentJSON, err := readCurrentAppJSON()
	if err != nil {
		return ChangeAnalysis{}, errors.Wrap(err, "reading current app.json")
	}

	currentDef, err := appdef.Parse(currentJSON)
	if err != nil {
		return ChangeAnalysis{}, errors.Wrap(err, "parsing current app.json")
	}

	// Get previous app.json from git.
	previousJSON, err := getFileFromGit(ctx, appdef.JsonFileName, baseRef)
	if err != nil {
		return ChangeAnalysis{}, errors.Wrap(err, "getting previous app.json from git")
	}

	// Parse previous app.json.
	previousDef, err := appdef.Parse(previousJSON)
	if err != nil {
		return ChangeAnalysis{}, errors.Wrap(err, "parsing previous app.json")
	}

	// Compare definitions.
	return compareDefinitions(currentDef, previousDef), nil
}

// compareDefinitions compares two app definitions and determines if Terraform is needed.
func compareDefinitions(current, previous *appdef.Definition) ChangeAnalysis {
	// Quick check: if JSON is identical, skip.
	currentNorm := normaliseForComparison(current)
	previousNorm := normaliseForComparison(previous)

	if jsonEqual(currentNorm, previousNorm) {
		return ChangeAnalysis{
			Skip:   true,
			Reason: "app.json unchanged",
		}
	}

	// Check non-env fields (infrastructure config).
	currentInfra := removeEnvFields(currentNorm)
	previousInfra := removeEnvFields(previousNorm)

	if !jsonEqual(currentInfra, previousInfra) {
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
func analyseEnvChanges(current, previous *appdef.Definition) []AppChange {
	changes := make([]AppChange, 0)

	// Create a map of previous apps for easier lookup.
	previousApps := make(map[string]*appdef.App)
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

		// Check if env vars changed.
		currentEnv := normaliseEnv(currentApp.Env)
		previousEnv := normaliseEnv(previousApp.Env)
		envChanged := !jsonEqual(currentEnv, previousEnv)

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

// Helper functions.

// getFileFromGit retrieves a file from a specific git reference.
func getFileFromGit(ctx context.Context, path, ref string) ([]byte, error) {
	cmd := executil.NewCommand("git", "show", ref+":"+path)
	result, err := executil.DefaultRunner().Run(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "executing git show")
	}
	return []byte(result.Output), nil
}

// readCurrentAppJSON reads the current app.json file from disk.
func readCurrentAppJSON() ([]byte, error) {
	data, err := os.ReadFile(appdef.JsonFileName)
	if err != nil {
		return nil, errors.Wrap(err, "reading app.json")
	}
	return data, nil
}

// normaliseForComparison converts a Definition to a normalised map structure.
func normaliseForComparison(def *appdef.Definition) map[string]interface{} {
	data, _ := json.Marshal(def)
	var norm map[string]interface{}
	_ = json.Unmarshal(data, &norm)
	return norm
}

// normaliseEnv converts an Environment to a normalised map for comparison.
func normaliseEnv(env appdef.Environment) map[string]interface{} {
	data, _ := json.Marshal(env)
	var norm map[string]interface{}
	_ = json.Unmarshal(data, &norm)
	return norm
}

// removeEnvFields creates a deep copy of the data and removes all env fields.
func removeEnvFields(data map[string]interface{}) map[string]interface{} {
	// Deep copy the map.
	copied := deepCopyMap(data)

	// Remove shared.env.
	if shared, ok := copied["shared"].(map[string]interface{}); ok {
		delete(shared, "env")
	}

	// Remove env from each app.
	if apps, ok := copied["apps"].([]interface{}); ok {
		for _, app := range apps {
			if appMap, ok := app.(map[string]interface{}); ok {
				delete(appMap, "env")
			}
		}
	}

	return copied
}

// deepCopyMap creates a deep copy of a map.
func deepCopyMap(src map[string]interface{}) map[string]interface{} {
	data, _ := json.Marshal(src)
	var dst map[string]interface{}
	_ = json.Unmarshal(data, &dst)
	return dst
}

// jsonEqual compares two map structures for equality by marshalling to JSON.
func jsonEqual(a, b map[string]interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}
