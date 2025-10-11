package secrets

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

func setupTestSecretFiles(t *testing.T, fs afero.Fs, baseDir string) {
	t.Helper()

	environments := []string{env.Development, env.Staging, env.Production}
	for _, e := range environments {
		path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")

		// Create directory
		dir := filepath.Dir(path)
		require.NoError(t, fs.MkdirAll(dir, 0755))

		// Write plaintext secret file with actual secret values
		content := `# ` + e + ` environment secrets
SECRET_KEY: "test_secret_value_` + e + `"
API_TOKEN: "token_123_` + e + `"
DATABASE_PASSWORD: "super_secret_password_` + e + `"
`
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0600))
	}
}

func setupSOPSConfig(t *testing.T, fs afero.Fs, baseDir string, agePublicKey string) {
	t.Helper()

	sopsConfig := `creation_rules:
  - path_regex: secrets/.*\.yaml$
    age: ` + agePublicKey + `
`
	configPath := filepath.Join(baseDir, "resources", ".sops.yaml")
	dir := filepath.Dir(configPath)
	require.NoError(t, fs.MkdirAll(dir, 0755))
	require.NoError(t, afero.WriteFile(fs, configPath, []byte(sopsConfig), 0644))
}
