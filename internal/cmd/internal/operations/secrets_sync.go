// cmd/internal/operations/secrets_sync.go
package operations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

// SecretsSync adds missing secret placeholders to SOPS files based on app.json.
// Only works with unencrypted files.
func SecretsSync(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()
	fs := input.FS

	// Extract all SOPS references from app.json.
	secretRefs := extractSOPSReferences(app)
	if len(secretRefs) == 0 {
		fmt.Println("No secrets with source: 'sops' found in app.json")
		return nil
	}

	// Group secrets by file
	secretsByFile := groupSecretsByFile(secretRefs)

	fmt.Println("Syncing secrets from app.json...\n")

	var (
		totalAdded     int
		totalSkipped   int
		totalEncrypted int
	)

	// Process each file
	for filePath, secrets := range secretsByFile {
		result := processSecretFile(fs, filePath, secrets)

		fmt.Printf("✓ %s\n", filePath)
		for _, msg := range result.messages {
			fmt.Printf("  %s\n", msg)
		}
		fmt.Println()

		totalAdded += result.added
		totalSkipped += result.skipped
		if result.encrypted {
			totalEncrypted++
		}
	}

	// Print summary
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Summary: %d secrets added, %d skipped", totalAdded, totalSkipped)
	if totalEncrypted > 0 {
		fmt.Printf(", %d encrypted files skipped", totalEncrypted)
	}
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// secretReference represents a secret extracted from app.json
type secretReference struct {
	Key      string   // The secret key name
	FilePath string   // Path to the SOPS file
	AppNames []string // Apps that use this secret
}

// syncResult tracks the results of syncing a single file
type syncResult struct {
	added     int
	skipped   int
	encrypted bool
	messages  []string
}

// extractSOPSReferences finds all secrets with source: "sops" in app.json
func extractSOPSReferences(app *appdef.Definition) []secretReference {
	var refs []secretReference
	keyToRef := make(map[string]*secretReference) // file:key -> ref

	// Helper to process environment variables
	processEnvVars := func(envVars appdef.EnvVar, appName string) {
		for key, val := range envVars {
			if val.Source != appdef.EnvSourceSOPS {
				continue
			}

			// Parse the path: "secrets/production.yaml:PAYLOAD_SECRET"
			sopsPath, err := val.ParseSOPSPath()
			if err != nil {
				fmt.Printf("⚠ Warning: Invalid SOPS path for %s: %s\n", key, val.Path)
				continue
			}

			refKey := fmt.Sprintf("%s:%s", sopsPath.File, sopsPath.Key)
			if existing, ok := keyToRef[refKey]; ok {
				existing.AppNames = append(existing.AppNames, appName)
			} else {
				ref := &secretReference{
					Key:      sopsPath.Key,
					FilePath: sopsPath.File,
					AppNames: []string{appName},
				}
				keyToRef[refKey] = ref
				refs = append(refs, *ref)
			}
		}
	}

	// Extract from shared env
	if app.Shared.Env.Dev != nil {
		processEnvVars(app.Shared.Env.Dev, "shared")
	}
	if app.Shared.Env.Staging != nil {
		processEnvVars(app.Shared.Env.Staging, "shared")
	}
	if app.Shared.Env.Production != nil {
		processEnvVars(app.Shared.Env.Production, "shared")
	}

	// Extract from each app's env
	for _, appItem := range app.Apps {
		if appItem.Env.Dev != nil {
			processEnvVars(appItem.Env.Dev, appItem.Name)
		}
		if appItem.Env.Staging != nil {
			processEnvVars(appItem.Env.Staging, appItem.Name)
		}
		if appItem.Env.Production != nil {
			processEnvVars(appItem.Env.Production, appItem.Name)
		}
	}

	return refs
}

// groupSecretsByFile organizes secrets by their target file
func groupSecretsByFile(refs []secretReference) map[string][]secretReference {
	grouped := make(map[string][]secretReference)
	for _, ref := range refs {
		grouped[ref.FilePath] = append(grouped[ref.FilePath], ref)
	}
	return grouped
}

// processSecretFile adds missing secrets to a single SOPS file
func processSecretFile(fs afero.Fs, filePath string, secrets []secretReference) syncResult {
	result := syncResult{messages: []string{}}

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		result.messages = append(result.messages, fmt.Sprintf("✗ Error checking file: %v", err))
		return result
	}

	if !exists {
		// Create the file
		if err := createSecretFile(fs, filePath, secrets); err != nil {
			result.messages = append(result.messages, fmt.Sprintf("✗ Error creating file: %v", err))
			return result
		}
		result.added = len(secrets)
		result.messages = append(result.messages, fmt.Sprintf("• Created file with %d secrets", len(secrets)))
		return result
	}

	// Read existing file
	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		result.messages = append(result.messages, fmt.Sprintf("✗ Error reading file: %v", err))
		return result
	}

	// Check if file is encrypted
	if isEncrypted(content) {
		result.encrypted = true
		result.messages = append(result.messages, "⚠ File is encrypted - decrypt first with:")
		result.messages = append(result.messages, fmt.Sprintf("  sops %s", filePath))
		return result
	}

	// Parse existing YAML
	var existing map[string]any
	if err := yaml.Unmarshal(content, &existing); err != nil {
		result.messages = append(result.messages, fmt.Sprintf("✗ Error parsing YAML: %v", err))
		return result
	}

	if existing == nil {
		existing = make(map[string]any)
	}

	// Add missing secrets
	var additions strings.Builder
	for _, secret := range secrets {
		if _, exists := existing[secret.Key]; exists {
			result.skipped++
			result.messages = append(result.messages, fmt.Sprintf("• Skipped %s (already exists)", secret.Key))
		} else {
			result.added++
			appList := strings.Join(secret.AppNames, ", ")
			result.messages = append(result.messages, fmt.Sprintf("• Added %s (used by: %s)", secret.Key, appList))

			// Add to additions
			if additions.Len() == 0 {
				additions.WriteString("\n# Added by webkit secrets sync\n")
			}
			additions.WriteString(fmt.Sprintf("# Used by: %s\n", appList))
			additions.WriteString(fmt.Sprintf("%s: \"REPLACE_ME_%s\"\n", secret.Key, strings.ToUpper(secret.Key)))
		}
	}

	// If we added anything, append to file
	if result.added > 0 {
		updatedContent := append(content, []byte(additions.String())...)
		if err := afero.WriteFile(fs, filePath, updatedContent, 0644); err != nil {
			result.messages = append(result.messages, fmt.Sprintf("✗ Error writing file: %v", err))
			return result
		}
	}

	return result
}

// createSecretFile creates a new secret file with the given secrets
func createSecretFile(fs afero.Fs, filePath string, secrets []secretReference) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	var content strings.Builder
	env := strings.TrimSuffix(filepath.Base(filePath), ".yaml")

	content.WriteString(fmt.Sprintf("# %s environment secrets\n", env))
	content.WriteString("# Generated by webkit secrets sync\n\n")

	for _, secret := range secrets {
		appList := strings.Join(secret.AppNames, ", ")
		content.WriteString(fmt.Sprintf("# Used by: %s\n", appList))
		content.WriteString(fmt.Sprintf("%s: \"REPLACE_ME_%s\"\n\n", secret.Key, strings.ToUpper(secret.Key)))
	}

	return afero.WriteFile(fs, filePath, []byte(content.String()), 0644)
}

// isEncrypted checks if a file is SOPS encrypted
func isEncrypted(content []byte) bool {
	// SOPS encrypted files contain "sops:" metadata
	return strings.Contains(string(content), "sops:")
}
