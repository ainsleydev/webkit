package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppType_String(t *testing.T) {
	t.Parallel()

	got := AppTypeGoLang.String()
	assert.Equal(t, "golang", got)
	assert.IsType(t, "", got)
}

func TestApp_Language(t *testing.T) {
	t.Parallel()

	tt := []struct {
		input AppType
		want  string
	}{
		{input: AppTypeGoLang, want: "go"},
		{input: AppTypePayload, want: "js"},
		{input: AppTypeSvelteKit, want: "js"},
	}

	for _, test := range tt {
		t.Run(test.input.String(), func(t *testing.T) {
			t.Parallel()
			a := App{Type: test.input}
			got := a.Language()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestDomainType_String(t *testing.T) {
	t.Parallel()

	got := DomainTypePrimary.String()
	assert.Equal(t, "primary", got)
	assert.IsType(t, "", got)
}

func TestApp_OrderedCommands(t *testing.T) {
	t.Parallel()

	app := &App{
		Name: "web",
		Type: AppTypeGoLang,
		Path: "./",
	}

	err := app.applyDefaults()
	require.NoError(t, err)

	commands := app.OrderedCommands()
	require.Len(t, commands, 4)

	t.Log("In Order")
	{
		assert.Equal(t, "format", commands[0].Name)
		assert.Equal(t, "lint", commands[1].Name)
		assert.Equal(t, "test", commands[2].Name)
		assert.Equal(t, "build", commands[3].Name)
	}

	t.Log("Check CMD is Populated")
	{
		assert.Equal(t, "gofmt -w .", commands[0].Cmd)
		assert.Equal(t, "golangci-lint run", commands[1].Cmd)
		assert.Equal(t, "go test ./...", commands[2].Cmd)
		assert.Equal(t, "go build main.go", commands[3].Cmd)
	}
}

func TestMergeAllEnvironments(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def  Definition
		want Environment
	}{
		"Shared only": {
			def: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
					},
				},
			},
			want: Environment{
				Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
			},
		},
		"App overrides shared": {
			def: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
					},
				},
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "app"}},
						},
					},
				},
			},
			want: Environment{
				Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "app"}},
			},
		},
		"Multiple apps, last wins": {
			def: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
					},
				},
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "first"}},
						},
					},
					{
						Name: "app2",
						Env: Environment{
							Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "second"}},
						},
					},
				},
			},
			want: Environment{
				Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "second"}},
			},
		},
		"App only, no shared": {
			def: Definition{
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev: EnvVar{"BAR": {Source: EnvSourceValue, Value: "val"}},
						},
					},
				},
			},
			want: Environment{
				Dev: EnvVar{"BAR": {Source: EnvSourceValue, Value: "val"}},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.def.MergeAllEnvironments()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMergeAppEnvironment(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def       Definition
		appName   string
		want      Environment
		wantFound bool
	}{
		"App exists, shared only": {
			def: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
					},
				},
				Apps: []App{
					{Name: "app1"},
				},
			},
			appName: "app1",
			want: Environment{
				Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
			},
			wantFound: true,
		},
		"App exists, overrides shared": {
			def: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
					},
				},
				Apps: []App{
					{
						Name: "app1",
						Env: Environment{
							Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "app"}},
						},
					},
				},
			},
			appName: "app1",
			want: Environment{
				Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "app"}},
			},
			wantFound: true,
		},
		"App not found": {
			def: Definition{
				Shared: Shared{
					Env: Environment{
						Dev: EnvVar{"FOO": {Source: EnvSourceValue, Value: "shared"}},
					},
				},
			},
			appName:   "nonexistent",
			want:      Environment{},
			wantFound: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got, found := test.def.MergeAppEnvironment(test.appName)
			assert.Equal(t, test.wantFound, found)
			assert.Equal(t, test.want, got)
		})
	}
}
