package secrets

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
)

// ValidationResult contains the results of validating secrets
type ValidationResult struct {
	Valid          bool
	MissingSecrets []MissingSecret
	OrphanedKeys   []OrphanedKey
	Files          []FileValidation
}

// MissingSecret represents a secret referenced in app.json but not found in SOPS files
type MissingSecret struct {
	Key         string
	Environment string
	AppName     string
	ExpectedIn  string // Which SOPS file it should be in
}

// OrphanedKey represents a key in SOPS files not referenced in app.json
type OrphanedKey struct {
	Key      string
	FilePath string
}

// FileValidation contains validation results for a single file
type FileValidation struct {
	FilePath     string
	Exists       bool
	IsEncrypted  bool
	MissingKeys  []string
	OrphanedKeys []string
	Error        error
}

// ValidateConfig holds configuration for secret validation
type ValidateConfig struct {
	FS               afero.Fs
	AppDef           *appdef.Definition
	CheckOrphans     bool // Whether to report keys in SOPS files not in app.json
	AllowEncrypted   bool // Whether encrypted files should be validated
	SecretsDirectory string
}

// Validate checks that all secrets referenced in app.json exist in their respective SOPS files
func Validate(cfg ValidateConfig) (*ValidationResult, error) {
	if cfg.SecretsDirectory == "" {
		cfg.SecretsDirectory = "FilePath"
	}

	result := &ValidationResult{
		Valid:          true,
		MissingSecrets: []MissingSecret{},
		OrphanedKeys:   []OrphanedKey{},
		Files:          []FileValidation{},
	}

	// Build map of all secrets that should exist
	expectedSecrets := buildExpectedSecretsMap(cfg.AppDef)

	// Validate each environment's secret file
	environments := []string{"development", "staging", "production"}

	for _, env := range environments {
		filePath := fmt.Sprintf("%s/%s.yaml", cfg.SecretsDirectory, env)
		fileResult := validateSecretFile(cfg.FS, filePath, expectedSecrets[env], cfg)

		result.Files = append(result.Files, fileResult)

		// Aggregate missing secrets
		for _, key := range fileResult.MissingKeys {
			result.MissingSecrets = append(result.MissingSecrets, MissingSecret{
				Key:         key,
				Environment: env,
				AppName:     findAppForSecret(cfg.AppDef, env, key),
				ExpectedIn:  filePath,
			})
			result.Valid = false
		}

		// Aggregate orphaned keys if checking
		if cfg.CheckOrphans {
			for _, key := range fileResult.OrphanedKeys {
				result.OrphanedKeys = append(result.OrphanedKeys, OrphanedKey{
					Key:      key,
					FilePath: filePath,
				})
			}
		}

		if fileResult.Error != nil {
			result.Valid = false
		}
	}

	return result, nil
}

// buildExpectedSecretsMap creates a map of environment -> set of expected secret keys
func buildExpectedSecretsMap(def *appdef.Definition) map[string]map[string]bool {
	expected := map[string]map[string]bool{
		"development": make(map[string]bool),
		"staging":     make(map[string]bool),
		"production":  make(map[string]bool),
	}

	// Collect secrets from shared environment
	//def.Shared.Env.Walk(func(env string, name string, value appdef.EnvValue) {
	//	if value.Source == appdef.EnvSourceSOPS {
	//		expected[env][name] = true
	//	}
	//})
	//
	//// Collect secrets from each app's environment
	//for _, app := range def.Apps {
	//	app.Env.Walk(func(env string, name string, value appdef.EnvValue) {
	//		if value.Source == appdef.EnvSourceSOPS {
	//			expected[env][name] = true
	//		}
	//	})
	//}

	return expected
}

// validateSecretFile validates a single SOPS file against expected secrets
func validateSecretFile(fs afero.Fs, filePath string, expectedKeys map[string]bool, cfg ValidateConfig) FileValidation {
	result := FileValidation{
		FilePath:     filePath,
		MissingKeys:  []string{},
		OrphanedKeys: []string{},
	}

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		result.Error = fmt.Errorf("checking file existence: %w", err)
		return result
	}
	result.Exists = exists

	if !exists {
		// All expected keys are missing
		for key := range expectedKeys {
			result.MissingKeys = append(result.MissingKeys, key)
		}
		return result
	}

	// Read file
	data, err := afero.ReadFile(fs, filePath)
	if err != nil {
		result.Error = fmt.Errorf("reading file: %w", err)
		return result
	}

	// Check if encrypted
	result.IsEncrypted = sops.IsContentEncrypted(data)

	if result.IsEncrypted && !cfg.AllowEncrypted {
		result.Error = fmt.Errorf("file is encrypted - cannot validate (decrypt first or use --allow-encrypted)")
		return result
	}

	// Parse YAML to get actual keys
	actualKeys, err := extractYAMLKeys(data)
	if err != nil {
		result.Error = fmt.Errorf("parsing YAML: %w", err)
		return result
	}

	// Find missing keys
	for expectedKey := range expectedKeys {
		if !actualKeys[expectedKey] {
			result.MissingKeys = append(result.MissingKeys, expectedKey)
		}
	}

	// Find orphaned keys (if checking)
	if cfg.CheckOrphans {
		for actualKey := range actualKeys {
			if !expectedKeys[actualKey] {
				result.OrphanedKeys = append(result.OrphanedKeys, actualKey)
			}
		}
	}

	return result
}

// findAppForSecret finds which app uses a given secret
func findAppForSecret(def *appdef.Definition, env string, key string) string {
	// Check shared first
	sharedEnv := getEnvForEnvironment(def.Shared.Env, env)
	if sharedEnv != nil {
		if val, exists := sharedEnv[key]; exists && val.Source == appdef.EnvSourceSOPS {
			return "shared"
		}
	}

	// Check each app
	for _, app := range def.Apps {
		appEnv := getEnvForEnvironment(app.Env, env)
		if appEnv != nil {
			if val, exists := appEnv[key]; exists && val.Source == appdef.EnvSourceSOPS {
				return app.Name
			}
		}
	}

	return "unknown"
}

// getEnvForEnvironment extracts the appropriate EnvVar for an environment name
func getEnvForEnvironment(env appdef.Environment, envName string) appdef.EnvVar {
	switch envName {
	case "development":
		return env.Dev
	case "staging":
		return env.Staging
	case "production":
		return env.Production
	default:
		return nil
	}
}

// extractYAMLKeys parses YAML and returns all top-level keys
func extractYAMLKeys(data []byte) (map[string]bool, error) {
	// Reuse existing YAML parsing logic
	parsed := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	keys := make(map[string]bool)
	for k := range parsed {
		// Skip SOPS metadata keys
		if !strings.HasPrefix(k, "sops_") && k != "sops" {
			keys[k] = true
		}
	}

	return keys, nil
}
