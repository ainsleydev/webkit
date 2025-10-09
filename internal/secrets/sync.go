package secrets

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
)

// SyncConfig configures the sync operation
type SyncConfig struct {
	FS     afero.Fs
	AppDef *appdef.Definition
}

// SecretInfo contains metadata about an added secret
type SecretInfo struct {
	Key      string
	AppNames []string
}

// Sync performs the secrets sync operation.
// It extracts SOPS references from app.json, groups them by file,
// and adds missing placeholders to each secret file.
func Sync(cfg SyncConfig) (*SyncResults, error) {
	def := cfg.AppDef

	// 1. Merge all the environment variables and gather them as
	// individual references.
	var refs []reference
	for _, app := range def.Apps {
		e, ok := def.MergeAppEnvironment(app.Name)
		if !ok {
			return &SyncResults{}, nil
		}
		e.Walk(func(env string, name string, val appdef.EnvValue) {
			if val.Source != appdef.EnvSourceSOPS {
				return
			}
			refs = append(refs, reference{
				Key:         name,
				Environment: env,
				AppNames:    []string{app.Name},
			})
		})
	}

	// 2. Deduplicate (same key in same env used by multiple apps).
	refs = deduplicateByKey(refs)
	if len(refs) == 0 {
		return &SyncResults{}, nil
	}

	// 3. Group by file (environment determines the file).
	grouped := make(map[string][]reference)
	for _, ref := range refs {
		grouped[ref.GetFilePath()] = append(grouped[ref.GetFilePath()], ref)
	}

	// 4. Process each file
	results := &SyncResults{Files: make([]SyncResult, 0, len(grouped))}
	for filePath, secrets := range grouped {
		results.Files = append(results.Files, processFile(cfg.FS, filePath, secrets))
	}

	return results, nil
}

type reference struct {
	Key         string   // e.g., "PAYLOAD_SECRET"
	Environment string   // e.g., "production", "staging", "development"
	AppNames    []string // Apps using this secret
}

// GetFilePath returns the SOPS file path for this reference
func (r reference) GetFilePath() string {
	return filepath.Join(FilePath, r.Environment+".yaml")
}

// DeduplicateByKey removes duplicate keys (same key used by multiple apps)
// and merges their app names
func deduplicateByKey(refs []reference) []reference {
	keyMap := make(map[string]*reference)

	for _, ref := range refs {
		key := fmt.Sprintf("%s:%s", ref.Environment, ref.Key)
		if existing, ok := keyMap[key]; ok {
			existing.AppNames = append(existing.AppNames, ref.AppNames...)
		} else {
			refCopy := ref
			keyMap[key] = &refCopy
		}
	}

	result := make([]reference, 0, len(keyMap))
	for _, ref := range keyMap {
		result = append(result, *ref)
	}

	return result
}

// processFile processes a single secret file by adding missing placeholders.
// It checks if the file exists, is encrypted, and adds any missing secret keys.
func processFile(fs afero.Fs, filePath string, secrets []reference) SyncResult {
	result := SyncResult{
		FilePath:     filePath,
		AddedSecrets: []SecretInfo{},
	}

	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		result.Error = fmt.Errorf("checking file existence: %w", err)
		return result
	}

	if !exists {
		result.WasMissing = true
		return result
	}

	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		result.Error = fmt.Errorf("reading file: %w", err)
		return result
	}

	if sops.IsContentEncrypted(content) {
		result.WasEncrypted = true
		return result
	}

	// Parse existing keys
	var data map[string]any
	if err := yaml.Unmarshal(content, &data); err != nil {
		result.Error = fmt.Errorf("parsing YAML: %w", err)
		return result
	}

	if data == nil {
		data = make(map[string]any)
	}

	// Process each secret reference and check if they exist
	// in the file; if they don't add a placeholder.
	var additions bytes.Buffer
	for _, secret := range secrets {
		if _, exists = data[secret.Key]; exists {
			result.Skipped++
			continue
		}

		result.Added++
		result.AddedSecrets = append(result.AddedSecrets, SecretInfo{
			Key:      secret.Key,
			AppNames: secret.AppNames,
		})

		additions.WriteString(fmt.Sprintf("%s: \"REPLACE_ME_%s\"\n",
			secret.Key, strings.ToUpper(secret.Key)))
	}

	// Write back to the file if any of the secrets
	// need scaffolding to the file.
	if result.Added > 0 {
		updatedContent := append(content, []byte(additions.String())...)
		if err = afero.WriteFile(fs, filePath, updatedContent, 0644); err != nil {
			result.Error = fmt.Errorf("writing file: %w", err)
			return result
		}
	}

	return result
}
