package appdef

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

// validDefinition returns a valid Definition for testing.
func validDefinition() *Definition {
	return &Definition{
		WebkitVersion: "1.0.0",
		Project: Project{
			Name:        "test-project",
			Title:       "Test Project",
			Description: "Test description",
			Repo:        GitHubRepo{Owner: "test", Name: "repo"},
		},
		Apps: []App{
			{
				Name:  "test-app",
				Title: "Test App",
				Type:  AppTypeGoLang,
				Path:  "/apps/test",
				Infra: Infra{
					Provider: ResourceProviderDigitalOcean,
					Type:     "vm",
				},
				Domains: []Domain{{Name: "test.example.com"}},
			},
		},
	}
}

func TestDefinition_Validate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input        *Definition
		setup        func(afero.Fs)
		wantErr      bool
		wantErrCount int
	}{
		"Valid Definition": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Domains = []Domain{{Name: "example.com"}}
				d.Resources = []Resource{{
					Name:     "db",
					Type:     ResourceTypePostgres,
					Provider: ResourceProviderDigitalOcean,
				}}
				return d
			}(),
			setup: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/test", 0o755))
			},
			wantErr: false,
		},
		"Domain With Protocol": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Domains = []Domain{{Name: "https://example.com"}}
				return d
			}(),
			setup: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/test", 0o755))
			},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Non-existent App Path": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Path = "/apps/nonexistent"
				return d
			}(),
			setup:        func(fs afero.Fs) {},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Terraform-managed VM Without Domains": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].TerraformManaged = ptr.BoolPtr(true)
				d.Apps[0].Domains = []Domain{}
				return d
			}(),
			setup: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/test", 0o755))
			},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Invalid Env Resource Reference": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Env = Environment{
					Production: EnvVar{
						"DATABASE_URL": EnvValue{
							Source: EnvSourceResource,
							Value:  "nonexistent.connection_url",
						},
					},
				}
				d.Resources = []Resource{{
					Name:     "db",
					Type:     ResourceTypePostgres,
					Provider: ResourceProviderDigitalOcean,
				}}
				return d
			}(),
			setup: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/test", 0o755))
			},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Multiple Validation Errors": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Path = "/apps/nonexistent"
				d.Apps[0].Domains = []Domain{{Name: "https://example.com"}}
				return d
			}(),
			setup:        func(fs afero.Fs) {},
			wantErr:      true,
			wantErrCount: 2, // path + domain
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			if test.setup != nil {
				test.setup(fs)
			}

			errs := test.input.Validate(fs)

			if test.wantErr {
				require.NotNil(t, errs)
				assert.Len(t, errs, test.wantErrCount)
			} else {
				assert.Nil(t, errs)
			}
		})
	}
}

func TestDefinition_ValidateDomains(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input        *Definition
		wantErr      bool
		wantErrCount int
	}{
		"Valid Domains": {
			input: &Definition{
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
			wantErr: false,
		},
		"Domain With HTTPS Protocol": {
			input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "https://example.com"},
						},
					},
				},
			},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Domain With HTTP Protocol": {
			input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "http://example.com"},
						},
					},
				},
			},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Multiple Apps With Protocol Errors": {
			input: &Definition{
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
			wantErr:      true,
			wantErrCount: 2,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.input.validateDomains()

			if test.wantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.wantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_ValidateAppPaths(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input        *Definition
		setup        func(afero.Fs)
		wantErr      bool
		wantErrCount int
	}{
		"Valid Paths": {
			input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/app1"},
					{Name: "app2", Path: "/apps/app2"},
				},
			},
			setup: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/app1", 0o755))
				require.NoError(t, fs.MkdirAll("/apps/app2", 0o755))
			},
			wantErr: false,
		},
		"Non-existent Path": {
			input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/nonexistent"},
				},
			},
			setup:        func(fs afero.Fs) {},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Empty Path Is Skipped": {
			input: &Definition{
				Apps: []App{
					{Name: "app1", Path: ""},
				},
			},
			setup:   func(fs afero.Fs) {},
			wantErr: false,
		},
		"Mixed Valid And Invalid Paths": {
			input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/app1"},
					{Name: "app2", Path: "/apps/nonexistent"},
				},
			},
			setup: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/app1", 0o755))
			},
			wantErr:      true,
			wantErrCount: 1,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			if test.setup != nil {
				test.setup(fs)
			}

			errs := test.input.validateAppPaths(fs)

			if test.wantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.wantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_ValidateTerraformManagedVMs(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input        *Definition
		wantErr      bool
		wantErrCount int
	}{
		"Terraform-managed VM With Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{{Name: "example.com"}},
					},
				},
			},
			wantErr: false,
		},
		"Terraform-managed VM Without Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Terraform-managed App Without Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "app"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			wantErr:      true,
			wantErrCount: 1,
		},
		"Non-terraform-managed VM Without Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(false),
						Domains:          []Domain{},
					},
				},
			},
			wantErr: false,
		},
		"Terraform-managed Non-VM Type Without Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "other"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			wantErr: false,
		},
		"Default Terraform-managed (Nil) VM Without Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: nil, // defaults to true
						Domains:          []Domain{},
					},
				},
			},
			wantErr:      true,
			wantErrCount: 1,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.input.validateTerraformManagedVMs()

			if test.wantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.wantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_ValidateEnvReferences(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input        *Definition
		wantErr      bool
		wantErrCount int
	}{
		"Valid Resource Reference": {
			input: &Definition{
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
			wantErr: false,
		},
		"Non-existent Resource": {
			input: &Definition{
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
			wantErr:      true,
			wantErrCount: 1,
		},
		"Invalid Output For Resource Type": {
			input: &Definition{
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
			wantErr:      true,
			wantErrCount: 1,
		},
		"Invalid Reference Format": {
			input: &Definition{
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
			wantErr:      true,
			wantErrCount: 1,
		},
		"Value Source Not Validated": {
			input: &Definition{
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
			wantErr: false,
		},
		"Shared Env With Valid Reference": {
			input: &Definition{
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
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.input.validateEnvReferences()

			if test.wantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.wantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
