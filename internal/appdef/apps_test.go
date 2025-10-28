package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
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

func TestApp_ShouldUseNPM(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		appType AppType
		usesNPM *bool
		want    bool
	}{
		"Payload Default":        {appType: AppTypePayload, usesNPM: nil, want: true},
		"SvelteKit Default":      {appType: AppTypeSvelteKit, usesNPM: nil, want: true},
		"GoLang Default":         {appType: AppTypeGoLang, usesNPM: nil, want: false},
		"Payload Explicit False": {appType: AppTypePayload, usesNPM: ptr.BoolPtr(false), want: false},
		"GoLang Explicit True":   {appType: AppTypeGoLang, usesNPM: ptr.BoolPtr(true), want: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := App{
				Type:    test.appType,
				UsesNPM: test.usesNPM,
			}
			got := a.ShouldUseNPM()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_OrderedCommands(t *testing.T) {
	t.Parallel()

	t.Run("Missing Skipped", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name:     "web",
			Type:     AppTypeGoLang,
			Path:     "./",
			Commands: map[Command]CommandSpec{},
		}

		commands := app.OrderedCommands()
		assert.Len(t, commands, 0)
	})

	t.Run("Default Populated", func(t *testing.T) {
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
	})
}

func TestApp_MergeEnvironments(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		app    App
		shared Environment
		want   Environment
	}{
		"No Shared Env": {
			app: App{
				Name: "app1",
				Env: Environment{
					Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
					Staging:    EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "app1"}},
					Production: EnvVar{"KEY3": EnvValue{Source: EnvSourceValue, Value: "app1"}},
				},
			},
			shared: Environment{},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
				Staging:    EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "app1"}},
				Production: EnvVar{"KEY3": EnvValue{Source: EnvSourceValue, Value: "app1"}},
			},
		},
		"App Overrides Shared": {
			app: App{
				Name: "app1",
				Env: Environment{
					Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}},
					Staging:    EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "app1"}},
					Production: EnvVar{"KEY3": EnvValue{Source: EnvSourceValue, Value: "app1"}},
				},
			},
			shared: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY3": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
			},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "app1"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Staging:    EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY2": EnvValue{Source: EnvSourceValue, Value: "app1"}, "KEY3": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Production: EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}, "KEY3": EnvValue{Source: EnvSourceValue, Value: "app1"}},
			},
		},
		"Empty App Env Uses Shared": {
			app: App{
				Name: "app1",
				Env:  Environment{},
			},
			shared: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Staging:    EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Production: EnvVar{"KEY3": EnvValue{Source: EnvSourceValue, Value: "shared"}},
			},
			want: Environment{
				Dev:        EnvVar{"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Staging:    EnvVar{"KEY2": EnvValue{Source: EnvSourceValue, Value: "shared"}},
				Production: EnvVar{"KEY3": EnvValue{Source: EnvSourceValue, Value: "shared"}},
			},
		},
		"Different Source Types Override": {
			app: App{
				Name: "app1",
				Env: Environment{
					Dev: EnvVar{
						"KEY1": EnvValue{Source: EnvSourceSOPS, Path: "secrets/app.yaml:KEY1"},
					},
				},
			},
			shared: Environment{
				Dev: EnvVar{
					"KEY1": EnvValue{Source: EnvSourceValue, Value: "shared_value"},
					"KEY2": EnvValue{Source: EnvSourceResource, Value: "shared.resource"},
				},
			},
			want: Environment{
				Dev: EnvVar{
					"KEY1": EnvValue{Source: EnvSourceSOPS, Path: "secrets/app.yaml:KEY1"},
					"KEY2": EnvValue{Source: EnvSourceResource, Value: "shared.resource"},
				},
				Staging:    EnvVar{},
				Production: EnvVar{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.app.MergeEnvironments(tc.shared)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestApp_IsTerraformManaged(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		terraformManaged *bool
		want             bool
	}{
		"Nil defaults to true":     {terraformManaged: nil, want: true},
		"Explicit false":            {terraformManaged: ptr.BoolPtr(false), want: false},
		"Explicit true":             {terraformManaged: ptr.BoolPtr(true), want: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			app := App{TerraformManaged: test.terraformManaged}
			got := app.IsTerraformManaged()
			assert.Equal(t, test.want, got)
		})
	}
}
