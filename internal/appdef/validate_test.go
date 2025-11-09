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
		WantErr      bool
		WantErrCount int
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
			WantErr: false,
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
			WantErr:      true,
			WantErrCount: 1,
		},
		"Non-existent App Path": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Path = "/apps/nonexistent"
				return d
			}(),
			setup:        func(fs afero.Fs) {},
			WantErr:      true,
			WantErrCount: 1,
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
			WantErr:      true,
			WantErrCount: 1,
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
			WantErr:      true,
			WantErrCount: 1,
		},
		"Multiple Validation Errors": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Path = "/apps/nonexistent"
				d.Apps[0].Domains = []Domain{{Name: "https://example.com"}}
				return d
			}(),
			setup:        func(fs afero.Fs) {},
			WantErr:      true,
			WantErrCount: 2, // path + domain
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

			if test.WantErr {
				require.NotNil(t, errs)
				assert.Len(t, errs, test.WantErrCount)
			} else {
				assert.Nil(t, errs)
			}
		})
	}
}

func TestDefinition_validateDomains(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		Input        *Definition
		WantErr      bool
		WantErrCount int
	}{
		"Valid Domains": {
			Input: &Definition{
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
			WantErr: false,
		},
		"Domain With HTTPS Protocol": {
			Input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "https://example.com"},
						},
					},
				},
			},
			WantErr:      true,
			WantErrCount: 1,
		},
		"Domain With HTTP Protocol": {
			Input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "http://example.com"},
						},
					},
				},
			},
			WantErr:      true,
			WantErrCount: 1,
		},
		"Multiple Apps With Protocol Errors": {
			Input: &Definition{
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
			WantErr:      true,
			WantErrCount: 2,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.Input.validateDomains()

			if test.WantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.WantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_validateAppPaths(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		Input        *Definition
		SetupFS      func(afero.Fs)
		WantErr      bool
		WantErrCount int
	}{
		"Valid Paths": {
			Input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/app1"},
					{Name: "app2", Path: "/apps/app2"},
				},
			},
			SetupFS: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/app1", 0o755))
				require.NoError(t, fs.MkdirAll("/apps/app2", 0o755))
			},
			WantErr: false,
		},
		"Non-existent Path": {
			Input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/nonexistent"},
				},
			},
			SetupFS:      func(fs afero.Fs) {},
			WantErr:      true,
			WantErrCount: 1,
		},
		"Empty Path Is Skipped": {
			Input: &Definition{
				Apps: []App{
					{Name: "app1", Path: ""},
				},
			},
			SetupFS: func(fs afero.Fs) {},
			WantErr: false,
		},
		"Mixed Valid And Invalid Paths": {
			Input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/app1"},
					{Name: "app2", Path: "/apps/nonexistent"},
				},
			},
			SetupFS: func(fs afero.Fs) {
				require.NoError(t, fs.MkdirAll("/apps/app1", 0o755))
			},
			WantErr:      true,
			WantErrCount: 1,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			if test.SetupFS != nil {
				test.SetupFS(fs)
			}

			errs := test.Input.validateAppPaths(fs)

			if test.WantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.WantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_validateTerraformManagedVMs(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		Input        *Definition
		WantErr      bool
		WantErrCount int
	}{
		"Terraform-managed VM With Domains": {
			Input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{{Name: "example.com"}},
					},
				},
			},
			WantErr: false,
		},
		"Terraform-managed VM Without Domains": {
			Input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			WantErr:      true,
			WantErrCount: 1,
		},
		"Terraform-managed App Without Domains": {
			Input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "app"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			WantErr:      true,
			WantErrCount: 1,
		},
		"Non-terraform-managed VM Without Domains": {
			Input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(false),
						Domains:          []Domain{},
					},
				},
			},
			WantErr: false,
		},
		"Terraform-managed Non-VM Type Without Domains": {
			Input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "other"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			WantErr: false,
		},
		"Default Terraform-managed (Nil) VM Without Domains": {
			Input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: nil, // defaults to true
						Domains:          []Domain{},
					},
				},
			},
			WantErr:      true,
			WantErrCount: 1,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.Input.validateTerraformManagedVMs()

			if test.WantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.WantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestDefinition_validateEnvReferences(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		Input        *Definition
		WantErr      bool
		WantErrCount int
	}{
		"Valid Resource Reference": {
			Input: &Definition{
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
			WantErr: false,
		},
		"Non-existent Resource": {
			Input: &Definition{
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
			WantErr:      true,
			WantErrCount: 1,
		},
		"Invalid Output For Resource Type": {
			Input: &Definition{
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
			WantErr:      true,
			WantErrCount: 1,
		},
		"Invalid Reference Format": {
			Input: &Definition{
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
			WantErr:      true,
			WantErrCount: 1,
		},
		"Value Source Not Validated": {
			Input: &Definition{
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
			WantErr: false,
		},
		"Shared Env With Valid Reference": {
			Input: &Definition{
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
			WantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.Input.validateEnvReferences()

			if test.WantErr {
				require.NotEmpty(t, errs)
				assert.Len(t, errs, test.WantErrCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
