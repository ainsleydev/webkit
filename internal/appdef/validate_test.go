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
		input    *Definition
		setup    func(afero.Fs)
		wantErrs []string
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
			wantErrs: []string{},
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
			wantErrs: []string{
				`app "test-app": domain "https://example.com" should not contain protocol prefix`,
			},
		},
		"Non-existent App Path": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Path = "/apps/nonexistent"
				return d
			}(),
			setup: func(fs afero.Fs) {},
			wantErrs: []string{
				`app "test-app": path "/apps/nonexistent" does not exist`,
			},
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
			wantErrs: []string{
				`app "test-app": terraform-managed VM/app must have at least one domain configured`,
			},
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
			wantErrs: []string{
				`env var "DATABASE_URL" in production references non-existent resource "nonexistent"`,
			},
		},
		"Multiple Validation Errors": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps[0].Path = "/apps/nonexistent"
				d.Apps[0].Domains = []Domain{{Name: "https://example.com"}}
				return d
			}(),
			setup: func(fs afero.Fs) {},
			wantErrs: []string{
				`domain "https://example.com" should not contain protocol prefix`,
				`path "/apps/nonexistent" does not exist`,
			},
		},
		"Empty Apps List": {
			input: &Definition{
				WebkitVersion: "1.0.0",
				Project: Project{
					Name:        "test-project",
					Title:       "Test Project",
					Description: "Test description",
					Repo:        GitHubRepo{Owner: "test", Name: "repo"},
				},
				Apps: []App{},
			},
			setup: func(fs afero.Fs) {},
			wantErrs: []string{
				"validation failed on 'min' tag",
			},
		},
		"Complex Multi-Error Scenario": {
			input: func() *Definition {
				d := validDefinition()
				d.Apps = append(d.Apps, App{
					Name:             "app2",
					Title:            "App 2",
					Type:             AppTypePayload,
					Path:             "/apps/missing",
					TerraformManaged: ptr.BoolPtr(true),
					Infra:            Infra{Type: "vm", Provider: ResourceProviderDigitalOcean},
					Domains:          []Domain{{Name: "http://app2.com"}},
				})
				d.Apps[0].Domains = []Domain{{Name: "ftp://example.com"}}
				d.Apps[0].Path = "/apps/also-missing"
				return d
			}(),
			setup: func(fs afero.Fs) {},
			wantErrs: []string{
				`domain "ftp://example.com" should not contain protocol prefix`,
				`domain "http://app2.com" should not contain protocol prefix`,
				`path "/apps/also-missing" does not exist`,
				`path "/apps/missing" does not exist`,
			},
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

			if len(test.wantErrs) == 0 {
				assert.Nil(t, errs, "expected no errors")
			} else {
				require.Len(t, errs, len(test.wantErrs), "unexpected number of errors")

				for i, wantErr := range test.wantErrs {
					assert.Contains(t, errs[i].Error(), wantErr,
						"error message should contain expected substring")
				}
			}
		})
	}
}

func TestDefinition_ValidateDomains(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input    *Definition
		wantErrs []string
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
			wantErrs: []string{},
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
			wantErrs: []string{
				`app "test-app": domain "https://example.com" should not contain protocol prefix`,
			},
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
			wantErrs: []string{
				`app "test-app": domain "http://example.com" should not contain protocol prefix`,
			},
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
			wantErrs: []string{
				`app "app1": domain "https://example.com"`,
				`app "app2": domain "http://api.example.com"`,
			},
		},
		"Mixed Valid And Invalid Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "valid.com"},
							{Name: "https://invalid.com"},
							{Name: "also-valid.com"},
						},
					},
				},
			},
			wantErrs: []string{
				`domain "https://invalid.com" should not contain protocol prefix`,
			},
		},
		"Domain With FTP Protocol": {
			input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Domains: []Domain{
							{Name: "ftp://files.example.com"},
						},
					},
				},
			},
			wantErrs: []string{
				`domain "ftp://files.example.com" should not contain protocol prefix`,
			},
		},
		"Empty Domains List": {
			input: &Definition{
				Apps: []App{
					{
						Name:    "test-app",
						Domains: []Domain{},
					},
				},
			},
			wantErrs: []string{},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.input.validateDomains()

			if len(test.wantErrs) == 0 {
				assert.Empty(t, errs, "expected no errors")
			} else {
				require.Len(t, errs, len(test.wantErrs), "unexpected number of errors")

				for i, wantErr := range test.wantErrs {
					assert.Contains(t, errs[i].Error(), wantErr,
						"error message should contain expected substring")
				}
			}
		})
	}
}

func TestDefinition_ValidateAppPaths(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input    *Definition
		setup    func(afero.Fs)
		wantErrs []string
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
			wantErrs: []string{},
		},
		"Non-existent Path": {
			input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/nonexistent"},
				},
			},
			setup: func(fs afero.Fs) {},
			wantErrs: []string{
				`app "app1": path "/apps/nonexistent" does not exist`,
			},
		},
		"Empty Path Is Skipped": {
			input: &Definition{
				Apps: []App{
					{Name: "app1", Path: ""},
				},
			},
			setup:    func(fs afero.Fs) {},
			wantErrs: []string{},
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
			wantErrs: []string{
				`app "app2": path "/apps/nonexistent" does not exist`,
			},
		},
		"Multiple Non-existent Paths": {
			input: &Definition{
				Apps: []App{
					{Name: "app1", Path: "/apps/missing1"},
					{Name: "app2", Path: "/apps/missing2"},
					{Name: "app3", Path: "/apps/missing3"},
				},
			},
			setup: func(fs afero.Fs) {},
			wantErrs: []string{
				`app "app1": path "/apps/missing1" does not exist`,
				`app "app2": path "/apps/missing2" does not exist`,
				`app "app3": path "/apps/missing3" does not exist`,
			},
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

			if len(test.wantErrs) == 0 {
				assert.Empty(t, errs, "expected no errors")
			} else {
				require.Len(t, errs, len(test.wantErrs), "unexpected number of errors")

				for i, wantErr := range test.wantErrs {
					assert.Contains(t, errs[i].Error(), wantErr,
						"error message should contain expected substring")
				}
			}
		})
	}
}

func TestDefinition_ValidateTerraformManagedVMs(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input    *Definition
		wantErrs []string
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
			wantErrs: []string{},
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
			wantErrs: []string{
				`app "test-app": terraform-managed VM/app must have at least one domain configured`,
			},
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
			wantErrs: []string{
				`app "test-app": terraform-managed VM/app must have at least one domain configured`,
			},
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
			wantErrs: []string{},
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
			wantErrs: []string{},
		},
		"Default Terraform-managed (Nil) VM Without Domains": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "test-app",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: nil,
						Domains:          []Domain{},
					},
				},
			},
			wantErrs: []string{
				`app "test-app": terraform-managed VM/app must have at least one domain configured`,
			},
		},
		"Multiple Apps With Mixed Configurations": {
			input: &Definition{
				Apps: []App{
					{
						Name:             "app1",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{{Name: "app1.com"}},
					},
					{
						Name:             "app2",
						Infra:            Infra{Type: "vm"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
					{
						Name:             "app3",
						Infra:            Infra{Type: "app"},
						TerraformManaged: ptr.BoolPtr(true),
						Domains:          []Domain{},
					},
				},
			},
			wantErrs: []string{
				`app "app2": terraform-managed VM/app must have at least one domain configured`,
				`app "app3": terraform-managed VM/app must have at least one domain configured`,
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.input.validateTerraformManagedVMs()

			if len(test.wantErrs) == 0 {
				assert.Empty(t, errs, "expected no errors")
			} else {
				require.Len(t, errs, len(test.wantErrs), "unexpected number of errors")

				for i, wantErr := range test.wantErrs {
					assert.Contains(t, errs[i].Error(), wantErr,
						"error message should contain expected substring")
				}
			}
		})
	}
}

func TestDefinition_ValidateEnvReferences(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input    *Definition
		wantErrs []string
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
			wantErrs: []string{},
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
			wantErrs: []string{
				`env var "DATABASE_URL" in production references non-existent resource "nonexistent"`,
			},
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
			wantErrs: []string{
				`env var "DATABASE_URL" in production references invalid output "invalid_output" for resource "db"`,
			},
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
			wantErrs: []string{
				`env var "DATABASE_URL" in production has invalid resource reference format "invalid-format"`,
			},
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
			wantErrs: []string{},
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
			wantErrs: []string{},
		},
		"Multiple Invalid References": {
			input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Env: Environment{
							Production: EnvVar{
								"DB_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "missing.connection_url",
								},
								"CACHE_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "cache.invalid_field",
								},
							},
						},
					},
				},
			},
			wantErrs: []string{
				`env var "DB_URL" in production references non-existent resource "missing"`,
				`env var "CACHE_URL" in production references non-existent resource "cache"`,
			},
		},
		"Dev And Staging Env References": {
			input: &Definition{
				Apps: []App{
					{
						Name: "test-app",
						Env: Environment{
							Dev: EnvVar{
								"DB_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "db.connection_url",
								},
							},
							Staging: EnvVar{
								"DB_URL": EnvValue{
									Source: EnvSourceResource,
									Value:  "missing.connection_url",
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
			wantErrs: []string{
				`env var "DB_URL" in staging references non-existent resource "missing"`,
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			errs := test.input.validateEnvReferences()

			if len(test.wantErrs) == 0 {
				assert.Empty(t, errs, "expected no errors")
			} else {
				require.Len(t, errs, len(test.wantErrs), "unexpected number of errors")

				for i, wantErr := range test.wantErrs {
					assert.Contains(t, errs[i].Error(), wantErr,
						"error message should contain expected substring")
				}
			}
		})
	}
}
