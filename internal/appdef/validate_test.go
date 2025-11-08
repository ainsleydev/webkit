package appdef

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinition_Validate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def         *Definition
		setupFS     func(afero.Fs)
		expectError bool
		errorCount  int
	}{
		"valid definition": {
			def: &Definition{
				Apps: []App{
					{
						Name:  "test-app",
						Path:  "/apps/test",
						Infra: Infra{Type: "vm"},
						Domains: []Domain{
							{Name: "example.com"},
						},
					},
				},
				Resources: []Resource{
					{
						Name:     "db",
						Type:     ResourceTypePostgres,
						Provider: ResourceProviderDigitalOcean,
					},
				},
			},
			setupFS: func(fs afero.Fs) {
				fs.MkdirAll("/apps/test", 0755)
			},
			expectError: false,
		},
		"domain with protocol": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Path: "/apps/test",
						Domains: []Domain{
							{Name: "https://example.com"},
						},
					},
				},
			},
			setupFS: func(fs afero.Fs) {
				fs.MkdirAll("/apps/test", 0755)
			},
			expectError: true,
			errorCount:  1,
		},
		"non-existent app path": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Path: "/apps/nonexistent",
					},
				},
			},
			setupFS:     func(fs afero.Fs) {},
			expectError: true,
			errorCount:  1,
		},
		"terraform-managed VM without domains": {
			def: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Path:             "/apps/test",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: boolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			setupFS: func(fs afero.Fs) {
				fs.MkdirAll("/apps/test", 0755)
			},
			expectError: true,
			errorCount:  1,
		},
		"invalid env resource reference": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Path: "/apps/test",
						Env: Environment{
							Production: EnvVar{
								"DATABASE_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "nonexistent.connection_url",
								},
							},
						},
					},
				},
				Resources: []Resource{
					{
						Name:     "db",
						Type:     ResourceTypePostgres,
						Provider: ResourceProviderDigitalOcean,
					},
				},
			},
			setupFS: func(fs afero.Fs) {
				fs.MkdirAll("/apps/test", 0755)
			},
			expectError: true,
			errorCount:  1,
		},
		"multiple validation errors": {
			def: &Definition{
				Apps: []App{
					{
						Name:  "test-app",
						Path:  "/apps/nonexistent",
						Infra: Infra{Type: "vm"},
						Domains: []Domain{
							{Name: "https://example.com"},
						},
					},
				},
			},
			setupFS:     func(fs afero.Fs) {},
			expectError: true,
			errorCount:  2, // path + domain
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			if test.setupFS != nil {
				test.setupFS(fs)
			}

			errs := test.def.Validate(fs)

			if test.expectError {
				require.NotNil(t, errs)
				assert.Len(t, errs, test.errorCount)
			} else {
				assert.Nil(t, errs)
			}
		})
	}
}

func TestDefinition_validateDomains(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def         *Definition
		expectError bool
		errorCount  int
	}{
		"valid domains": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "example.com"},
							{Name: "api.example.com"},
						},
					},
				},
			},
			expectError: false,
		},
		"domain with https protocol": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "https://example.com"},
						},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		"domain with http protocol": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "http://example.com"},
						},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		"multiple apps with protocol errors": {
			def: &Definition{
				Apps: []App{
					{
						Name: "app1",
						Domains: []Domain{
							{Name: "https://example.com"},
						},
					},
					{
						Name: "app2",
						Domains: []Domain{
							{Name: "http://api.example.com"},
						},
					},
				},
			},
			expectError: true,
			errorCount:  2,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.def.validateDomains()

			if test.expectError {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.errorCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_validateAppPaths(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def         *Definition
		setupFS     func(afero.Fs)
		expectError bool
		errorCount  int
	}{
		"valid paths": {
			def: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/app1"},
					{Name: "app2", Path: "/apps/app2"},
				},
			},
			setupFS: func(fs afero.Fs) {
				fs.MkdirAll("/apps/app1", 0755)
				fs.MkdirAll("/apps/app2", 0755)
			},
			expectError: false,
		},
		"non-existent path": {
			def: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/nonexistent"},
				},
			},
			setupFS:     func(fs afero.Fs) {},
			expectError: true,
			errorCount:  1,
		},
		"empty path is skipped": {
			def: &Definition{
				Apps: []App{
					{Name: "app1", Path: ""},
				},
			},
			setupFS:     func(fs afero.Fs) {},
			expectError: false,
		},
		"mixed valid and invalid paths": {
			def: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/app1"},
					{Name: "app2", Path: "/apps/nonexistent"},
				},
			},
			setupFS: func(fs afero.Fs) {
				fs.MkdirAll("/apps/app1", 0755)
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			if test.setupFS != nil {
				test.setupFS(fs)
			}

			errs := test.def.validateAppPaths(fs)

			if test.expectError {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.errorCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_validateTerraformManagedVMs(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def         *Definition
		expectError bool
		errorCount  int
	}{
		"terraform-managed VM with domains": {
			def: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: boolPtr(true),
						Domains:          []Domain{{Name: "example.com"}},
					},
				},
			},
			expectError: false,
		},
		"terraform-managed VM without domains": {
			def: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: boolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		"terraform-managed app without domains": {
			def: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "app"},
						TerraformManaged: boolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		"non-terraform-managed VM without domains": {
			def: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: boolPtr(false),
						Domains:          []Domain{},
					},
				},
			},
			expectError: false,
		},
		"terraform-managed non-VM type without domains": {
			def: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "other"},
						TerraformManaged: boolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			expectError: false,
		},
		"default terraform-managed (nil) VM without domains": {
			def: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: nil, // defaults to true
						Domains:          []Domain{},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.def.validateTerraformManagedVMs()

			if test.expectError {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.errorCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_validateEnvReferences(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def         *Definition
		expectError bool
		errorCount  int
	}{
		"valid resource reference": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Env: Environment{
							Production: EnvVar{
								"DATABASE_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "db.connection_url",
								},
							},
						},
					},
				},
				Resources: []Resource{
					{
						Name:     "db",
						Type:     ResourceTypePostgres,
						Provider: ResourceProviderDigitalOcean,
					},
				},
			},
			expectError: false,
		},
		"non-existent resource": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Env: Environment{
							Production: EnvVar{
								"DATABASE_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "nonexistent.connection_url",
								},
							},
						},
					},
				},
				Resources: []Resource{
					{
						Name:     "db",
						Type:     ResourceTypePostgres,
						Provider: ResourceProviderDigitalOcean,
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		"invalid output for resource type": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Env: Environment{
							Production: EnvVar{
								"DATABASE_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "db.invalid_output",
								},
							},
						},
					},
				},
				Resources: []Resource{
					{
						Name:     "db",
						Type:     ResourceTypePostgres,
						Provider: ResourceProviderDigitalOcean,
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		"invalid reference format": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Env: Environment{
							Production: EnvVar{
								"DATABASE_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "invalid-format",
								},
							},
						},
					},
				},
				Resources: []Resource{
					{
						Name:     "db",
						Type:     ResourceTypePostgres,
						Provider: ResourceProviderDigitalOcean,
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		"value source not validated": {
			def: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Env: Environment{
							Production: EnvVar{
								"API_URL": EnvValue{
									Source: EnvSourceValue,
									Value:  "https://api.example.com",
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		"shared env with valid reference": {
			def: &Definition{
				Shared: Shared{
					Env: Environment{
						Production: EnvVar{
							"S3_BUCKET": EnvValue{
								Source: EnvSourceResource,
								Value:  "storage.bucket_name",
							},
						},
					},
				},
				Apps: []App{
					{Name: "test-app"},
				},
				Resources: []Resource{
					{
						Name:     "storage",
						Type:     ResourceTypeS3,
						Provider: ResourceProviderDigitalOcean,
					},
				},
			},
			expectError: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.def.validateEnvReferences()

			if test.expectError {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.errorCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

// Helper function for tests
func boolPtr(b bool) *bool {
	return &b
}
