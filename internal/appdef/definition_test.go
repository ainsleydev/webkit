package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubLabels(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Definition
		want  []string
	}{
		"No Apps": {
			input: Definition{Apps: nil},
			want:  []string{"webkit"},
		},
		"Single App": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeSvelteKit},
				},
			},
			want: []string{"webkit", AppTypeSvelteKit.String()},
		},
		"Multiple Apps": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeGoLang},
					{Type: AppTypePayload},
				},
			},
			want: []string{
				"webkit",
				AppTypeGoLang.String(),
				AppTypePayload.String(),
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.GithubLabels()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestContainsGo(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Definition
		want  bool
	}{
		"Truthy": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeGoLang},
				},
			},
			want: true,
		},
		"Falsey": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeSvelteKit},
				},
			},
			want: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.ContainsGo()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestContainsJS(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Definition
		want  bool
	}{
		"Truthy": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeSvelteKit},
				},
			},
			want: true,
		},
		"Falsey": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeGoLang},
				},
			},
			want: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.ContainsJS()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMergeAllEnvironments(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Definition
		want  Environment
	}{
		"Empty Definition": {
			input: Definition{},
			want: Environment{
				Dev:        EnvVar{},
				Staging:    EnvVar{},
				Production: EnvVar{},
			},
		},
		"Only Shared Environment": {
			input: Definition{
				Shared: Shared{
					Env: Environment{
						Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_dev"}},
						Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_staging"}},
						Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_prod"}},
					},
				},
			},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_dev"}},
				Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_staging"}},
				Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_prod"}},
			},
		},
		"Single App Overrides Shared": {
			input: Definition{
				Shared: Shared{
					Env: Environment{
						Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}},
						Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
						Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
					},
				},
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
							Staging:    EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "app1"}},
							Production: EnvVar{"KEY3": EnvValue{Source: EnvSourceValue, Value: "app1"}},
						},
					},
				},
			},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "app1"}},
				Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY3": EnvValue{Source: EnvSourceValue, Value: "app1"}},
			},
		},
		"Multiple Apps Last Wins": {
			input: Definition{
				Shared: Shared{
					Env: Environment{
						Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}},
						Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
						Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
					},
				},
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
							Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
							Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
						},
					},
					{
						Name: "app2",
						Env: Environment{
							Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app2"}, "KEY3": EnvValue{Source: EnvSourceValue, Value: "app2"}},
							Staging:    EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "app2"}},
							Production: EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "app2"}},
						},
					},
				},
			},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app2"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY3": EnvValue{Source: EnvSourceValue, Value: "app2"}},
				Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "app2"}},
				Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "app2"}},
			},
		},
		"Apps With Empty Environments": {
			input: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
					},
				},
				Apps: []App{
					{Name: "app1", Env: Environment{}},
					{Name: "app2", Env: Environment{}},
				},
			},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Staging:    EnvVar{},
				Production: EnvVar{},
			},
		},
		"No Shared Only Apps": {
			input: Definition{
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
						},
					},
					{
						Name: "app2",
						Env: Environment{
							Dev: EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "app2"}},
						},
					},
				},
			},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "app2"}},
				Staging:    EnvVar{},
				Production: EnvVar{},
			},
		},
		"Different Source Types": {
			input: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{
							"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_value"},
							"KEY2": EnvValue{Source: EnvSourceResource, Value: "db.connection"},
						},
					},
				},
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev: EnvVar{
								"KEY1": EnvValue{Source: EnvSourceSOPS, Path: "secrets/prod.yaml:API_KEY"},
								"KEY3": EnvValue{Source: EnvSourceResource, Value: "cache.url"},
							},
						},
					},
				},
			},
			want: Environment{
				Dev: EnvVar{
					"KEY1": EnvValue{Source: EnvSourceSOPS, Path: "secrets/prod.yaml:API_KEY"},
					"KEY2": EnvValue{Source: EnvSourceResource, Value: "db.connection"},
					"KEY3": EnvValue{Source: EnvSourceResource, Value: "cache.url"},
				},
				Staging:    EnvVar{},
				Production: EnvVar{},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.input.MergeAllEnvironments()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestDefinition_FilterTerraformManaged(t *testing.T) {
	t.Parallel()

	trueVal := true
	falseVal := false

	tt := map[string]struct {
		input            Definition
		wantAppsCount    int
		wantResCount     int
		wantSkippedApps  []string
		wantSkippedRes   []string
	}{
		"Empty definition": {
			input: Definition{
				Project: Project{Name: "test-project"},
			},
			wantAppsCount:   0,
			wantResCount:    0,
			wantSkippedApps: []string{},
			wantSkippedRes:  []string{},
		},
		"All managed nil default": {
			input: Definition{
				Project: Project{Name: "test-project"},
				Apps: []App{
					{Name: "app1", TerraformManaged: nil},
					{Name: "app2", TerraformManaged: nil},
				},
				Resources: []Resource{
					{Name: "db", TerraformManaged: nil},
					{Name: "cache", TerraformManaged: nil},
				},
			},
			wantAppsCount:   2,
			wantResCount:    2,
			wantSkippedApps: []string{},
			wantSkippedRes:  []string{},
		},
		"All managed explicit true": {
			input: Definition{
				Project: Project{Name: "test-project"},
				Apps: []App{
					{Name: "app1", TerraformManaged: &trueVal},
					{Name: "app2", TerraformManaged: &trueVal},
				},
				Resources: []Resource{
					{Name: "db", TerraformManaged: &trueVal},
					{Name: "cache", TerraformManaged: &trueVal},
				},
			},
			wantAppsCount:   2,
			wantResCount:    2,
			wantSkippedApps: []string{},
			wantSkippedRes:  []string{},
		},
		"All unmanaged": {
			input: Definition{
				Project: Project{Name: "test-project"},
				Apps: []App{
					{Name: "app1", TerraformManaged: &falseVal},
					{Name: "app2", TerraformManaged: &falseVal},
				},
				Resources: []Resource{
					{Name: "db", TerraformManaged: &falseVal},
					{Name: "cache", TerraformManaged: &falseVal},
				},
			},
			wantAppsCount:   0,
			wantResCount:    0,
			wantSkippedApps: []string{"app1", "app2"},
			wantSkippedRes:  []string{"db", "cache"},
		},
		"Mixed managed and unmanaged apps": {
			input: Definition{
				Project: Project{Name: "test-project"},
				Apps: []App{
					{Name: "app1", TerraformManaged: nil},
					{Name: "app2", TerraformManaged: &falseVal},
					{Name: "app3", TerraformManaged: &trueVal},
				},
			},
			wantAppsCount:   2,
			wantResCount:    0,
			wantSkippedApps: []string{"app2"},
			wantSkippedRes:  []string{},
		},
		"Mixed managed and unmanaged resources": {
			input: Definition{
				Project: Project{Name: "test-project"},
				Resources: []Resource{
					{Name: "db", TerraformManaged: nil},
					{Name: "cache", TerraformManaged: &falseVal},
					{Name: "storage", TerraformManaged: &trueVal},
				},
			},
			wantAppsCount:   0,
			wantResCount:    2,
			wantSkippedApps: []string{},
			wantSkippedRes:  []string{"cache"},
		},
		"Complex mixed scenario": {
			input: Definition{
				Project: Project{Name: "test-project"},
				Apps: []App{
					{Name: "frontend", TerraformManaged: &trueVal},
					{Name: "backend", TerraformManaged: nil},
					{Name: "worker", TerraformManaged: &falseVal},
				},
				Resources: []Resource{
					{Name: "db", TerraformManaged: &trueVal},
					{Name: "cache", TerraformManaged: &falseVal},
					{Name: "queue", TerraformManaged: nil},
				},
				Shared: Shared{
					Env: Environment{
						Production: EnvVar{
							"KEY": EnvValue{Source: EnvSourceValue, Value: "value"},
						},
					},
				},
			},
			wantAppsCount:   2,
			wantResCount:    2,
			wantSkippedApps: []string{"worker"},
			wantSkippedRes:  []string{"cache"},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filtered, skipped := test.input.FilterTerraformManaged()

			// Check filtered apps count
			assert.Len(t, filtered.Apps, test.wantAppsCount, "filtered apps count mismatch")

			// Check filtered resources count
			assert.Len(t, filtered.Resources, test.wantResCount, "filtered resources count mismatch")

			// Check skipped apps
			assert.ElementsMatch(t, test.wantSkippedApps, skipped.Apps, "skipped apps mismatch")

			// Check skipped resources
			assert.ElementsMatch(t, test.wantSkippedRes, skipped.Resources, "skipped resources mismatch")

			// Verify that the filtered definition preserves other fields
			assert.Equal(t, test.input.Project, filtered.Project, "project should be preserved")
			assert.Equal(t, test.input.Shared, filtered.Shared, "shared config should be preserved")
			assert.Equal(t, test.input.WebkitVersion, filtered.WebkitVersion, "webkit version should be preserved")
		})
	}
}
