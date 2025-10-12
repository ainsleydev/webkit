package secrets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestSync(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		cfg := SyncConfig{
			FS:     fs,
			AppDef: &appdef.Definition{Apps: []appdef.App{}},
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Empty(t, results)
		assert.Equal(t, 0, results.TotalAdded())
	})

	t.Run("No SOPS References", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeGoLang,
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Empty(t, results.Files)
	})

	t.Run("Single App Single Secret", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		// Create secret file
		secretPath := FilePath + "/production.yaml"
		err := fs.MkdirAll("FilePath", 0755)
		require.NoError(t, err)
		err = afero.WriteFile(fs, secretPath, []byte("EXISTING_KEY: value\n"), 0644)
		require.NoError(t, err)

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"PAYLOAD_SECRET": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:PAYLOAD_SECRET",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Len(t, results.Files, 1)
		assert.Equal(t, 1, results.TotalAdded())
		assert.Equal(t, 0, results.TotalSkipped())

		// Verify file was updated
		content, err := afero.ReadFile(fs, secretPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "PAYLOAD_SECRET: \"REPLACE_ME_PAYLOAD_SECRET\"")
		assert.Contains(t, string(content), "EXISTING_KEY: value")
	})

	t.Run("Multiple Apps Same Secret", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		secretPath := filepath.Join(FilePath, "production.yaml")
		err := fs.MkdirAll("FilePath", 0755)
		require.NoError(t, err)
		err = afero.WriteFile(fs, secretPath, []byte(""), 0644)
		require.NoError(t, err)

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"DATABASE_URL": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:DATABASE_URL",
							},
						},
					},
				},
				{
					Name: "web",
					Type: appdef.AppTypeGoLang,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"DATABASE_URL": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:DATABASE_URL",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Len(t, results.Files, 1)
		assert.Equal(t, 1, results.TotalAdded())

		// Verify both apps are listed
		result := results.Files[0]
		assert.Len(t, result.AddedSecrets, 1)
		assert.ElementsMatch(t, []string{"cms", "web"}, result.AddedSecrets[0].AppNames)
	})

	t.Run("Multiple Environments", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(FilePath, 0755)
		require.NoError(t, err)

		// Create files for each environment
		for _, env := range []string{"development", "staging", "production"} {
			path := filepath.Join(FilePath, env+".yaml")
			err = afero.WriteFile(fs, path, []byte(""), 0644)
			require.NoError(t, err)
		}

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Env: appdef.Environment{
						Dev: appdef.EnvVar{
							"DEV_SECRET": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/development.yaml:DEV_SECRET",
							},
						},
						Staging: appdef.EnvVar{
							"STAGING_SECRET": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/staging.yaml:STAGING_SECRET",
							},
						},
						Production: appdef.EnvVar{
							"PROD_SECRET": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:PROD_SECRET",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Len(t, results.Files, 3)
		assert.Equal(t, 3, results.TotalAdded())
	})

	t.Run("Skip Existing Keys", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		secretPath := filepath.Join(FilePath, "production.yaml")
		err := fs.MkdirAll(FilePath, os.ModePerm)
		require.NoError(t, err)

		// Pre-populate with existing secret
		err = afero.WriteFile(fs, secretPath, []byte("API_KEY: existing_value\n"), 0644)
		require.NoError(t, err)

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeGoLang,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"API_KEY": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:API_KEY",
							},
							"NEW_KEY": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:NEW_KEY",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Equal(t, 1, results.TotalAdded())
		assert.Equal(t, 1, results.TotalSkipped())

		content, err := afero.ReadFile(fs, secretPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "API_KEY: existing_value")
		assert.Contains(t, string(content), "NEW_KEY: \"REPLACE_ME_NEW_KEY\"")
	})

	t.Run("File Missing", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"SECRET": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:SECRET",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Len(t, results.Files, 1)
		assert.True(t, results.Files[0].WasMissing)
		assert.Equal(t, 1, results.MissingCount())
	})

	t.Run("File Encrypted", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		secretPath := filepath.Join(FilePath, "production.yaml")
		err := fs.MkdirAll(FilePath, os.ModePerm)
		require.NoError(t, err)

		// Write encrypted content
		encryptedContent := []byte("sops:\n  kms: encrypted_data\nENC[AES256_GCM,data:...]")
		err = afero.WriteFile(fs, secretPath, encryptedContent, 0644)
		require.NoError(t, err)

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"SECRET": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:SECRET",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Len(t, results.Files, 1)
		assert.True(t, results.Files[0].WasEncrypted)
		assert.Equal(t, 1, results.EncryptedCount())
	})

	t.Run("Mixed Value and SOPS Sources", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		secretPath := filepath.Join(FilePath, "production.yaml")
		err := fs.MkdirAll(FilePath, 0755)
		require.NoError(t, err)
		err = afero.WriteFile(fs, secretPath, []byte(""), 0644)
		require.NoError(t, err)

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeGoLang,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"PUBLIC_KEY": {
								Source: appdef.EnvSourceValue,
								Value:  "public_value",
							},
							"SECRET_KEY": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:SECRET_KEY",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.Equal(t, 1, results.TotalAdded())

		// Only SOPS secret should be added
		result := results.Files[0]
		assert.Len(t, result.AddedSecrets, 1)
		assert.Equal(t, "SECRET_KEY", result.AddedSecrets[0].Key)
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		secretPath := filepath.Join(FilePath, "production.yaml")
		err := fs.MkdirAll(FilePath, os.ModePerm)
		require.NoError(t, err)

		// Write invalid YAML
		err = afero.WriteFile(fs, secretPath, []byte("invalid: yaml: content: bad"), 0644)
		require.NoError(t, err)

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"SECRET": {
								Source: appdef.EnvSourceSOPS,
								Path:   "secrets/production.yaml:SECRET",
							},
						},
					},
				},
			},
		}
		require.NoError(t, def.ApplyDefaults())

		cfg := SyncConfig{
			FS:     fs,
			AppDef: def,
		}

		results, err := Sync(cfg)
		require.NoError(t, err)
		assert.True(t, results.HasErrors())
		assert.NotNil(t, results.Files[0].Error)
	})
}

func TestDeduplicateByKey(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input []reference
		want  []reference
	}{
		"Empty Input": {
			input: []reference{},
			want:  []reference{},
		},
		"No Duplicates": {
			input: []reference{
				{Key: "KEY1", Environment: "production", AppNames: []string{"app1"}},
				{Key: "KEY2", Environment: "production", AppNames: []string{"app2"}},
			},
			want: []reference{
				{Key: "KEY1", Environment: "production", AppNames: []string{"app1"}},
				{Key: "KEY2", Environment: "production", AppNames: []string{"app2"}},
			},
		},
		"Same Key Different Environments": {
			input: []reference{
				{Key: "API_KEY", Environment: "staging", AppNames: []string{"app1"}},
				{Key: "API_KEY", Environment: "production", AppNames: []string{"app1"}},
			},
			want: []reference{
				{Key: "API_KEY", Environment: "staging", AppNames: []string{"app1"}},
				{Key: "API_KEY", Environment: "production", AppNames: []string{"app1"}},
			},
		},
		"Merge Apps for Same Key": {
			input: []reference{
				{Key: "DATABASE_URL", Environment: "production", AppNames: []string{"cms"}},
				{Key: "DATABASE_URL", Environment: "production", AppNames: []string{"web"}},
				{Key: "DATABASE_URL", Environment: "production", AppNames: []string{"api"}},
			},
			want: []reference{
				{Key: "DATABASE_URL", Environment: "production", AppNames: []string{"cms", "web", "api"}},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := deduplicateByKey(test.input)

			assert.Len(t, got, len(test.want))
			for i := range got {
				assert.Equal(t, test.want[i].Key, got[i].Key)
				assert.Equal(t, test.want[i].Environment, got[i].Environment)
				assert.ElementsMatch(t, test.want[i].AppNames, got[i].AppNames)
			}
		})
	}
}

func TestReference_GetFilePath(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		environment string
		want        string
	}{
		"Production": {
			environment: env.Production,
			want:        filepath.Join(FilePath, env.Production+".yaml"),
		},
		"Staging": {
			environment: env.Staging,
			want:        filepath.Join(FilePath, env.Staging+".yaml"),
		},
		"Development": {
			environment: env.Development,
			want:        filepath.Join(FilePath, env.Development+".yaml"),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ref := reference{
				Key:         "TEST_KEY",
				Environment: test.environment,
				AppNames:    []string{"test"},
			}

			got := ref.GetFilePath()
			assert.Equal(t, test.want, got)
		})
	}
}
