package secrets

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
)

// ResolvedEnvVar represents an environment variable that has been
// fully resolved (decrypted if necessary) and is ready for Terraform.
type ResolvedEnvVar struct {
	Key   string
	Value string
	Type  EnvVarType
}

// EnvVarType indicates how Terraform should handle the variable
type EnvVarType string

const (
	EnvVarTypeGeneral EnvVarType = "GENERAL"
	EnvVarTypeSecret  EnvVarType = "SECRET"
)

// Resolver decrypts and resolves environment variables from the app definition
type Resolver struct {
	sopsClient sops.EncrypterDecrypter
}

// NewResolver creates a new environment variable resolver with SOPS support
func NewResolver() (*Resolver, error) {
	// Initialize age provider (checks SOPS_AGE_KEY env var or ~/.config/webkit/age.key)
	ageProvider, err := age.NewProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize age provider: %w", err)
	}

	return &Resolver{
		sopsClient: sops.NewClient(ageProvider),
	}, nil
}

// ResolveEnvironmentVariables resolves all environment variables for a given environment.
// It handles three source types:
// - "value": static string value (pass through)
// - "resource": Terraform resource reference (formatted as "resource:name.output")
// - "sops": encrypted secret (decrypted using SOPS)
func (r *Resolver) ResolveEnvironmentVariables(envVars appdef.EnvVar) ([]ResolvedEnvVar, error) {
	resolved := make([]ResolvedEnvVar, 0, len(envVars))

	for key, config := range envVars {
		var resolvedVar ResolvedEnvVar
		var err error

		switch config.Source {
		case appdef.EnvSourceValue:
			// Static value - use as-is
			resolvedVar = ResolvedEnvVar{
				Key:   key,
				Value: config.Value,
				Type:  EnvVarTypeGeneral,
			}

		case appdef.EnvSourceResource:
			// Resource reference - format for Terraform to resolve later
			// e.g., "db.connection_url" → "resource:db.connection_url"
			resolvedVar = ResolvedEnvVar{
				Key:   key,
				Value: fmt.Sprintf("resource:%s", config.Value),
				Type:  EnvVarTypeGeneral,
			}

		case appdef.EnvSourceSOPS:
			// SOPS secret - decrypt now
			resolvedVar, err = r.decryptSOPSSecret(key, config.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt %s: %w", key, err)
			}

		default:
			return nil, fmt.Errorf("unknown env var source type: %s for key: %s", config.Source, key)
		}

		resolved = append(resolved, resolvedVar)
	}

	return resolved, nil
}

// decryptSOPSSecret decrypts a secret from a SOPS file and extracts the specified key.
// Path format: "secrets/production.yaml:PAYLOAD_SECRET"
func (r *Resolver) decryptSOPSSecret(envKey, path string) (ResolvedEnvVar, error) {
	// Parse SOPS path format: "file:key"
	parts := strings.SplitN(path, ":", 2)
	if len(parts) != 2 {
		return ResolvedEnvVar{}, fmt.Errorf("invalid sops path format: %s (expected 'file:key')", path)
	}

	filePath := parts[0]
	secretKey := parts[1]

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ResolvedEnvVar{}, fmt.Errorf("sops file not found: %s", filePath)
	}

	// Read encrypted file
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return ResolvedEnvVar{}, fmt.Errorf("failed to read sops file %s: %w", filePath, err)
	}

	// Check if file is actually encrypted
	if !sops.IsContentEncrypted(encryptedData) {
		return ResolvedEnvVar{}, fmt.Errorf("file is not encrypted: %s", filePath)
	}

	// Create temporary file for in-place decryption
	tmpFile, err := os.CreateTemp("", "webkit-sops-*.yaml")
	if err != nil {
		return ResolvedEnvVar{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write encrypted data to temp file
	if _, err := tmpFile.Write(encryptedData); err != nil {
		return ResolvedEnvVar{}, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Decrypt in-place using SOPS CLI
	if err := r.sopsClient.Decrypt(tmpFile.Name()); err != nil {
		return ResolvedEnvVar{}, fmt.Errorf("failed to decrypt sops file %s: %w", filePath, err)
	}

	// Read decrypted content
	decryptedData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return ResolvedEnvVar{}, fmt.Errorf("failed to read decrypted file: %w", err)
	}

	// Parse YAML to extract the specific key
	var data map[string]interface{}
	if err := yaml.Unmarshal(decryptedData, &data); err != nil {
		return ResolvedEnvVar{}, fmt.Errorf("failed to parse decrypted yaml: %w", err)
	}

	// Extract the value
	value, exists := data[secretKey]
	if !exists {
		return ResolvedEnvVar{}, fmt.Errorf("key %s not found in %s", secretKey, filePath)
	}

	// Convert to string
	valueStr, ok := value.(string)
	if !ok {
		return ResolvedEnvVar{}, fmt.Errorf("value for key %s is not a string", secretKey)
	}

	return ResolvedEnvVar{
		Key:   envKey,
		Value: valueStr,
		Type:  EnvVarTypeSecret, // Mark as secret for secure handling in Terraform
	}, nil
}

// MergeEnvVars merges shared and app-specific environment variables.
// App-specific values override shared values for the same key.
func MergeEnvVars(shared, appSpecific []ResolvedEnvVar) []ResolvedEnvVar {
	// Build map for quick lookup
	sharedMap := make(map[string]ResolvedEnvVar)
	for _, v := range shared {
		sharedMap[v.Key] = v
	}

	// Override with app-specific values
	for _, v := range appSpecific {
		sharedMap[v.Key] = v
	}

	// Convert back to slice
	merged := make([]ResolvedEnvVar, 0, len(sharedMap))
	for _, v := range sharedMap {
		merged = append(merged, v)
	}

	return merged
}
