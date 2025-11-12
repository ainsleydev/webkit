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
						"KEY1": EnvValue{Source: EnvSourceSOPS, Value: "KEY1"},
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
					"KEY1": EnvValue{Source: EnvSourceSOPS, Value: "KEY1"},
					"KEY2": EnvValue{Source: EnvSourceResource, Value: "shared.resource"},
				},
				Staging:    EnvVar{},
				Production: EnvVar{},
			},
		},
		"Shared Default Only": {
			app: App{
				Name: "app1",
				Env:  Environment{},
			},
			shared: Environment{
				Default: EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
			},
			want: Environment{
				Dev:        EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
				Staging:    EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
				Production: EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
			},
		},
		"App Default Overrides Shared Default": {
			app: App{
				Name: "app1",
				Env: Environment{
					Default: EnvVar{"API_KEY": EnvValue{Source: EnvSourceValue, Value: "app-key"}},
				},
			},
			shared: Environment{
				Default: EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
			},
			want: Environment{
				Dev:        EnvVar{"API_KEY": EnvValue{Source: EnvSourceValue, Value: "app-key"}},
				Staging:    EnvVar{"API_KEY": EnvValue{Source: EnvSourceValue, Value: "app-key"}},
				Production: EnvVar{"API_KEY": EnvValue{Source: EnvSourceValue, Value: "app-key"}},
			},
		},
		"Shared Default With App Specific Override": {
			app: App{
				Name: "app1",
				Env: Environment{
					Production: EnvVar{"API_KEY": EnvValue{Source: EnvSourceValue, Value: "prod-key"}},
				},
			},
			shared: Environment{
				Default: EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
			},
			want: Environment{
				Dev:        EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
				Staging:    EnvVar{"API_KEY": EnvValue{Source: EnvSourceSOPS}},
				Production: EnvVar{"API_KEY": EnvValue{Source: EnvSourceValue, Value: "prod-key"}},
			},
		},
		"Complex Default Merging": {
			app: App{
				Name: "app1",
				Env: Environment{
					Default:    EnvVar{"VAR1": EnvValue{Source: EnvSourceValue, Value: "app-default"}},
					Dev:        EnvVar{"VAR2": EnvValue{Source: EnvSourceValue, Value: "dev-specific"}},
					Production: EnvVar{"VAR1": EnvValue{Source: EnvSourceValue, Value: "prod-override"}},
				},
			},
			shared: Environment{
				Default: EnvVar{
					"VAR1": EnvValue{Source: EnvSourceSOPS},
					"VAR3": EnvValue{Source: EnvSourceValue, Value: "shared-default"},
				},
				Staging: EnvVar{"VAR4": EnvValue{Source: EnvSourceValue, Value: "staging-shared"}},
			},
			want: Environment{
				Dev: EnvVar{
					"VAR1": EnvValue{Source: EnvSourceValue, Value: "app-default"},
					"VAR2": EnvValue{Source: EnvSourceValue, Value: "dev-specific"},
					"VAR3": EnvValue{Source: EnvSourceValue, Value: "shared-default"},
				},
				Staging: EnvVar{
					"VAR1": EnvValue{Source: EnvSourceValue, Value: "app-default"},
					"VAR3": EnvValue{Source: EnvSourceValue, Value: "shared-default"},
					"VAR4": EnvValue{Source: EnvSourceValue, Value: "staging-shared"},
				},
				Production: EnvVar{
					"VAR1": EnvValue{Source: EnvSourceValue, Value: "prod-override"},
					"VAR3": EnvValue{Source: EnvSourceValue, Value: "shared-default"},
				},
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
		"Nil defaults to true": {terraformManaged: nil, want: true},
		"Explicit false":       {terraformManaged: ptr.BoolPtr(false), want: false},
		"Explicit true":        {terraformManaged: ptr.BoolPtr(true), want: true},
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

func TestApp_ShouldRelease(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		release *bool
		want    bool
	}{
		"Nil defaults to true": {release: nil, want: true},
		"Explicit false":       {release: ptr.BoolPtr(false), want: false},
		"Explicit true":        {release: ptr.BoolPtr(true), want: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			app := App{Build: Build{Release: test.release}}
			got := app.ShouldRelease()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_DefaultPort(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		appType AppType
		want    int
	}{
		"Payload":   {appType: AppTypePayload, want: 3000},
		"SvelteKit": {appType: AppTypeSvelteKit, want: 3001},
		"GoLang":    {appType: AppTypeGoLang, want: 8080},
		"Default":   {appType: "", want: 3000},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			app := App{Type: test.appType}
			got := app.defaultPort()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_ApplyDefaults_Port(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		app  App
		want int
	}{
		"Default SvelteKit app gets 3001": {
			app:  App{Name: "web", Type: AppTypeSvelteKit, Path: "./"},
			want: 3001,
		},
		"Explicit port is preserved": {
			app:  App{Name: "web", Type: AppTypeSvelteKit, Path: "./", Build: Build{Port: 4000}},
			want: 4000,
		},
		"Payload app gets 3000": {
			app:  App{Name: "cms", Type: AppTypePayload, Path: "./"},
			want: 3000,
		},
		"SvelteKit app gets 3001": {
			app:  App{Name: "web", Type: AppTypeSvelteKit, Path: "./"},
			want: 3001,
		},
		"GoLang app gets 8080": {
			app:  App{Name: "api", Type: AppTypeGoLang, Path: "./"},
			want: 8080,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			app := test.app

			err := app.applyDefaults()
			require.NoError(t, err)
			assert.Equal(t, test.want, app.Build.Port)
		})
	}
}

func TestApp_PrimaryDomain(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		domains []Domain
		want    string
	}{
		"Returns primary domain when present": {
			domains: []Domain{
				{Name: "example.com", Type: DomainTypePrimary, Zone: "example.com"},
				{Name: "www.example.com", Type: DomainTypeAlias, Zone: "example.com"},
			},
			want: "example.com",
		},
		"Returns first domain when no primary": {
			domains: []Domain{
				{Name: "www.example.com", Type: DomainTypeAlias, Zone: "example.com"},
				{Name: "example.com", Type: DomainTypeAlias, Zone: "example.com"},
			},
			want: "www.example.com",
		},
		"Returns primary domain even if not first": {
			domains: []Domain{
				{Name: "www.example.com", Type: DomainTypeAlias, Zone: "example.com"},
				{Name: "example.com", Type: DomainTypePrimary, Zone: "example.com"},
			},
			want: "example.com",
		},
		"Returns empty string when no domains": {
			domains: []Domain{},
			want:    "",
		},
		"Returns empty string when domains is nil": {
			domains: nil,
			want:    "",
		},
		"Returns first primary when multiple primaries": {
			domains: []Domain{
				{Name: "example.com", Type: DomainTypePrimary, Zone: "example.com"},
				{Name: "example.org", Type: DomainTypePrimary, Zone: "example.org"},
			},
			want: "example.com",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			app := App{Domains: test.domains}
			got := app.PrimaryDomain()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_ResolvedTools(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		app  App
		want map[string]string
	}{
		"GoLang defaults": {
			app: App{Type: AppTypeGoLang},
			want: map[string]string{
				"golangci-lint": "latest",
				"templ":         "latest",
				"sqlc":          "latest",
			},
		},
		"Payload no defaults": {
			app:  App{Type: AppTypePayload},
			want: map[string]string{},
		},
		"SvelteKit no defaults": {
			app:  App{Type: AppTypeSvelteKit},
			want: map[string]string{},
		},
		"Custom override": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"templ": "v0.2.543",
				},
			},
			want: map[string]string{
				"golangci-lint": "latest",
				"templ":         "v0.2.543",
				"sqlc":          "latest",
			},
		},
		"Add custom tool": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"buf": "v1.28.1",
				},
			},
			want: map[string]string{
				"golangci-lint": "latest",
				"templ":         "latest",
				"sqlc":          "latest",
				"buf":           "v1.28.1",
			},
		},
		"Disable default tool with empty string": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"sqlc": "",
				},
			},
			want: map[string]string{
				"golangci-lint": "latest",
				"templ":         "latest",
			},
		},
		"Disable default tool with disabled": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"templ": "disabled",
				},
			},
			want: map[string]string{
				"golangci-lint": "latest",
				"sqlc":          "latest",
			},
		},
		"Multiple overrides": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"golangci-lint": "v1.55.2",
					"templ":         "disabled",
					"buf":           "v1.28.1",
				},
			},
			want: map[string]string{
				"golangci-lint": "v1.55.2",
				"sqlc":          "latest",
				"buf":           "v1.28.1",
			},
		},
		"Only custom tools": {
			app: App{
				Type: AppTypePayload,
				Tools: map[string]string{
					"custom-tool": "v1.0.0",
				},
			},
			want: map[string]string{
				"custom-tool": "v1.0.0",
			},
		},
		"Nil tools map": {
			app: App{
				Type:  AppTypeGoLang,
				Tools: nil,
			},
			want: map[string]string{
				"golangci-lint": "latest",
				"templ":         "latest",
				"sqlc":          "latest",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.app.ResolvedTools()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_InstallCommands(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		app  App
		want []string
	}{
		"GoLang defaults": {
			app: App{Type: AppTypeGoLang},
			want: []string{
				"go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
				"go install github.com/a-h/templ/cmd/templ@latest",
				"go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
			},
		},
		"No tools": {
			app:  App{Type: AppTypePayload},
			want: []string{},
		},
		"Custom version": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"templ": "v0.2.543",
					"sqlc":  "disabled",
				},
			},
			want: []string{
				"go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
				"go install github.com/a-h/templ/cmd/templ@v0.2.543",
			},
		},
		"Full install path": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"golangci-lint":                   "disabled",
					"templ":                           "disabled",
					"sqlc":                            "disabled",
					"github.com/custom/tool/cmd/tool": "v1.0.0",
				},
			},
			want: []string{
				"go install github.com/custom/tool/cmd/tool@v1.0.0",
			},
		},
		"Mixed known and custom": {
			app: App{
				Type: AppTypeGoLang,
				Tools: map[string]string{
					"buf":                      "v1.28.1",
					"github.com/custom/mytool": "v2.0.0",
				},
			},
			want: []string{
				"go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
				"go install github.com/a-h/templ/cmd/templ@latest",
				"go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
				"go install github.com/bufbuild/buf/cmd/buf@v1.28.1",
				"go install github.com/custom/mytool@v2.0.0",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.app.InstallCommands()
			// Sort both slices to ensure consistent comparison,
			// since map iteration order is not guaranteed.
			assert.ElementsMatch(t, test.want, got)
		})
	}
}
