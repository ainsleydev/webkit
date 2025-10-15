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
