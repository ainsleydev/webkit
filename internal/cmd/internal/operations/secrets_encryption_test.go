package operations

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
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

func TestSecretsEncrypt_Integration(t *testing.T) {
	// Skip in short mode since this is an integration test
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	t.Run("Encrypts All Environment Files", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		baseDir := "test-project"

		// Generate a real age key pair for this test
		// Note: In real integration tests, you might use a fixture key
		// For now, this will use the actual age.NewProvider() which reads from env/config
		setupTestSecretFiles(t, fs, baseDir)

		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
			BaseDir:     baseDir,
		}

		err := SecretsEncrypt(t.Context(), input)
		require.NoError(t, err)

		t.Log("All files are encrypted with SOPS")
		{
			environments := []string{env.Development, env.Staging, env.Production}
			for _, e := range environments {
				path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")

				exists, err := afero.Exists(fs, path)
				require.NoError(t, err)
				assert.True(t, exists, "File should exist: "+path)

				content, err := afero.ReadFile(fs, path)
				require.NoError(t, err)

				fmt.Print(string(content))

				// Verify SOPS encryption markers are present
				assert.Contains(t, string(content), "sops:", "File should contain SOPS metadata")
				assert.Contains(t, string(content), "mac:", "File should contain SOPS MAC")
				assert.Contains(t, string(content), "lastmodified:", "File should contain SOPS timestamp")

				// Verify original plaintext is NOT visible
				assert.NotContains(t, string(content), "test_secret_value_"+e, "Plaintext secrets should be encrypted")
			}
		}
	})

	t.Run("Handles Missing Files Gracefully", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		baseDir := "test-project"

		// Don't create any secret files
		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
			BaseDir:     baseDir,
		}

		// Should not panic or return error, just logs
		err := SecretsEncrypt(t.Context(), input)
		assert.NoError(t, err, "Should handle missing files gracefully")
	})

	t.Run("Preserves File Structure", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		baseDir := "test-project"

		setupTestSecretFiles(t, fs, baseDir)

		// Create additional files to ensure we don't affect them
		otherFile := filepath.Join(baseDir, secrets.FilePath, "README.md")
		require.NoError(t, afero.WriteFile(fs, otherFile, []byte("# Secrets"), 0644))

		err := SecretsEncrypt(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
			BaseDir:     baseDir,
		})
		require.NoError(t, err)

		t.Log("Non-secret files are unchanged")
		{
			content, err := afero.ReadFile(fs, otherFile)
			require.NoError(t, err)
			assert.Equal(t, "# Secrets", string(content), "Non-secret files should not be modified")
		}
	})
}

func TestSecretsDecrypt_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	t.Run("Decrypts All Environment Files", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		baseDir := "test-project"

		setupTestSecretFiles(t, fs, baseDir)

		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
			BaseDir:     baseDir,
		}

		// First encrypt
		err := SecretsEncrypt(t.Context(), input)
		require.NoError(t, err)

		// Then decrypt
		err = SecretsDecrypt(t.Context(), input)
		require.NoError(t, err)

		t.Log("All files are decrypted")
		{
			environments := []string{env.Development, env.Staging, env.Production}
			for _, e := range environments {
				path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")

				content, err := afero.ReadFile(fs, path)
				require.NoError(t, err)

				// Verify SOPS metadata is gone
				assert.NotContains(t, string(content), "sops:", "SOPS metadata should be removed")
				assert.NotContains(t, string(content), "mac:", "SOPS MAC should be removed")

				// Verify original content is restored
				assert.Contains(t, string(content), "SECRET_KEY", "Should contain secret keys")
				assert.Contains(t, string(content), "test_secret_value_"+e, "Should contain decrypted values")
			}
		}
	})

	t.Run("Handles Already Decrypted Files", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		baseDir := "test-project"

		// Create plaintext files
		setupTestSecretFiles(t, fs, baseDir)

		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
			BaseDir:     baseDir,
		}

		// Try to decrypt already plaintext files - should handle gracefully
		err := SecretsDecrypt(t.Context(), input)
		assert.NoError(t, err, "Should handle plaintext files gracefully")

		t.Log("Plaintext files remain unchanged")
		{
			path := filepath.Join(baseDir, secrets.FilePath, env.Development+".yaml")
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)
			assert.Contains(t, string(content), "SECRET_KEY", "Content should be preserved")
		}
	})
}

func TestSecretsEncryptDecrypt_RoundTrip_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	fs := afero.NewMemMapFs()
	baseDir := "test-project"

	setupTestSecretFiles(t, fs, baseDir)

	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: &appdef.Definition{},
		BaseDir:     baseDir,
	}

	// Store original content
	originalContent := make(map[string]string)
	environments := []string{env.Development, env.Staging, env.Production}
	for _, e := range environments {
		path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")
		content, err := afero.ReadFile(fs, path)
		require.NoError(t, err)
		originalContent[e] = string(content)
	}

	t.Log("Step 1: Encrypt all files")
	{
		err := SecretsEncrypt(t.Context(), input)
		require.NoError(t, err)

		for _, e := range environments {
			path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)

			encryptedContent := string(content)
			assert.Contains(t, encryptedContent, "sops:", "File should be encrypted")
			assert.NotEqual(t, originalContent[e], encryptedContent, "Content should change after encryption")
		}
	}

	t.Log("Step 2: Decrypt all files")
	{
		err := SecretsDecrypt(t.Context(), input)
		require.NoError(t, err)

		for _, e := range environments {
			path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)

			decryptedContent := string(content)
			assert.NotContains(t, decryptedContent, "sops:", "SOPS metadata should be removed")
			assert.Equal(t, originalContent[e], decryptedContent, "Decrypted content should match original exactly")
		}
	}

	t.Log("Step 3: Re-encrypt to verify idempotency")
	{
		err := SecretsEncrypt(t.Context(), input)
		require.NoError(t, err)

		for _, e := range environments {
			path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)
			assert.Contains(t, string(content), "sops:", "File should be encrypted again")
		}
	}

	t.Log("Step 4: Decrypt again to verify full round-trip")
	{
		err := SecretsDecrypt(t.Context(), input)
		require.NoError(t, err)

		for _, e := range environments {
			path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)
			assert.Equal(t, originalContent[e], string(content), "Multiple round-trips should preserve content")
		}
	}
}

func TestSecretsEncrypt_FilePermissions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	fs := afero.NewMemMapFs()
	baseDir := "test-project"

	setupTestSecretFiles(t, fs, baseDir)

	err := SecretsEncrypt(t.Context(), cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: &appdef.Definition{},
		BaseDir:     baseDir,
	})
	require.NoError(t, err)

	t.Log("Encrypted files maintain secure permissions")
	{
		environments := []string{env.Development, env.Staging, env.Production}
		for _, e := range environments {
			path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")

			info, err := fs.Stat(path)
			require.NoError(t, err)

			mode := info.Mode().Perm()
			// Secret files should be readable only by owner (0600 or 0400)
			assert.True(t, mode&0077 == 0,
				"File %s should not be readable by group or others, got: %o", path, mode)
		}
	}
}

func TestSecretsEncryptDecrypt_MultipleEnvironments_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	fs := afero.NewMemMapFs()
	baseDir := "test-project"

	setupTestSecretFiles(t, fs, baseDir)

	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: &appdef.Definition{},
		BaseDir:     baseDir,
	}

	// Encrypt all
	err := SecretsEncrypt(t.Context(), input)
	require.NoError(t, err)

	t.Log("Each environment has unique encrypted content")
	{
		devPath := filepath.Join(baseDir, secrets.FilePath, env.Development+".yaml")
		stagingPath := filepath.Join(baseDir, secrets.FilePath, env.Staging+".yaml")
		prodPath := filepath.Join(baseDir, secrets.FilePath, env.Production+".yaml")

		devContent, err := afero.ReadFile(fs, devPath)
		require.NoError(t, err)
		stagingContent, err := afero.ReadFile(fs, stagingPath)
		require.NoError(t, err)
		prodContent, err := afero.ReadFile(fs, prodPath)
		require.NoError(t, err)

		// All should be encrypted
		assert.Contains(t, string(devContent), "sops:")
		assert.Contains(t, string(stagingContent), "sops:")
		assert.Contains(t, string(prodContent), "sops:")

		// But should have different encrypted values
		assert.NotEqual(t, string(devContent), string(stagingContent),
			"Dev and staging should have different encrypted content")
		assert.NotEqual(t, string(stagingContent), string(prodContent),
			"Staging and production should have different encrypted content")
	}

	// Decrypt all
	err = SecretsDecrypt(t.Context(), input)
	require.NoError(t, err)

	t.Log("Each environment preserves unique plaintext content")
	{
		environments := []string{env.Development, env.Staging, env.Production}
		for _, e := range environments {
			path := filepath.Join(baseDir, secrets.FilePath, e+".yaml")
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)

			assert.Contains(t, string(content), "test_secret_value_"+e,
				"Environment %s should have its unique secret value", e)
		}
	}
}
