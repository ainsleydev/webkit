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

		t.Log("Check Tools are Populated")
		{
			require.NotNil(t, app.Tools)
			assert.Equal(t, "latest", app.Tools["golangci-lint"].Version)
			assert.Equal(t, "latest", app.Tools["templ"].Version)
			assert.Equal(t, "latest", app.Tools["sqlc"].Version)
			assert.Len(t, app.Tools, 3)
		}
	})

	t.Run("User Tools Preserved", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name: "api",
			Type: AppTypeGoLang,
			Path: "./",
			Tools: map[string]Tool{
				"templ": {Type: "go", Name: "github.com/a-h/templ/cmd/templ", Version: "v0.2.543"},
				"buf":   {Type: "go", Name: "github.com/bufbuild/buf/cmd/buf", Version: "v1.28.1"},
			},
		}

		err := app.applyDefaults()
		require.NoError(t, err)

		t.Log("Check User Tools are Preserved")
		{
			require.NotNil(t, app.Tools)
			assert.Equal(t, "v0.2.543", app.Tools["templ"].Version)
			assert.Equal(t, "v1.28.1", app.Tools["buf"].Version)
		}

		t.Log("Check Default Tools are Added")
		{
			assert.Equal(t, "latest", app.Tools["golangci-lint"].Version)
			assert.Equal(t, "latest", app.Tools["sqlc"].Version)
		}

		t.Log("Check All Tools Present")
		{
			assert.Len(t, app.Tools, 4)
		}
	})

	t.Run("Payload Apps Have No Default Tools", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name: "cms",
			Type: AppTypePayload,
			Path: "./",
		}

		err := app.applyDefaults()
		require.NoError(t, err)

		t.Log("Check No Default Tools for Payload")
		{
			require.NotNil(t, app.Tools)
			assert.Len(t, app.Tools, 0)
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

func TestApp_InstallCommands(t *testing.T) {
	t.Parallel()

	t.Run("GoLang defaults", func(t *testing.T) {
		t.Parallel()

		app := App{Type: AppTypeGoLang}
		err := app.applyDefaults()
		require.NoError(t, err)

		got := app.InstallCommands()
		assert.Len(t, got, 3)
		assert.Contains(t, got, "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
		assert.Contains(t, got, "go install github.com/a-h/templ/cmd/templ@latest")
		assert.Contains(t, got, "go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest")
	})

	t.Run("Custom Go tool with version", func(t *testing.T) {
		t.Parallel()

		app := App{
			Type: AppTypeGoLang,
			Tools: map[string]Tool{
				"templ": {Type: "go", Name: "github.com/a-h/templ/cmd/templ", Version: "v0.2.543"},
			},
		}
		err := app.applyDefaults()
		require.NoError(t, err)

		got := app.InstallCommands()
		assert.Len(t, got, 3)
		assert.Contains(t, got, "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
		assert.Contains(t, got, "go install github.com/a-h/templ/cmd/templ@v0.2.543")
		assert.Contains(t, got, "go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest")
	})

	t.Run("pnpm tool", func(t *testing.T) {
		t.Parallel()

		app := App{
			Type: AppTypeGoLang,
			Tools: map[string]Tool{
				"eslint": {Type: "pnpm", Name: "eslint", Version: "8.0.0"},
			},
		}
		err := app.applyDefaults()
		require.NoError(t, err)

		got := app.InstallCommands()
		assert.Contains(t, got, "pnpm add -g eslint@8.0.0")
	})

	t.Run("Custom install command via script type", func(t *testing.T) {
		t.Parallel()

		app := App{
			Type: AppTypeGoLang,
			Tools: map[string]Tool{
				"custom": {
					Type:    "script",
					Install: "curl -sSL https://example.com/install.sh | sh",
				},
			},
		}
		err := app.applyDefaults()
		require.NoError(t, err)

		got := app.InstallCommands()
		assert.Contains(t, got, "curl -sSL https://example.com/install.sh | sh")
	})

	t.Run("Install override for any type", func(t *testing.T) {
		t.Parallel()

		app := App{
			Type: AppTypeGoLang,
			Tools: map[string]Tool{
				"custom": {
					Type:    "go",
					Name:    "github.com/foo/bar",
					Version: "v1.0.0",
					Install: "custom install command",
				},
			},
		}
		err := app.applyDefaults()
		require.NoError(t, err)

		got := app.InstallCommands()
		assert.Contains(t, got, "custom install command")
	})
}

func TestApp_CommandOrderPreservation(t *testing.T) {
	t.Parallel()

	t.Run("Order preserved through marshal/unmarshal", func(t *testing.T) {
		t.Parallel()

		// Create JSON with custom command order
		jsonData := []byte(`{
			"name": "api",
			"title": "API",
			"type": "golang",
			"path": "./api",
			"build": {
				"dockerfile": "Dockerfile",
				"port": 8080
			},
			"infra": {
				"provider": "digitalocean",
				"type": "vm"
			},
			"commands": {
				"generate": {"command": "TEMPL_EXPERIMENT=rawgo go generate ./..."},
				"build": {"command": "go build main.go"},
				"format": {"command": "go fmt ./..."},
				"lint": {"command": "echo"},
				"test": {"command": "go test ./..."}
			}
		}`)

		// Unmarshal
		var app App
		err := json.Unmarshal(jsonData, &app)
		require.NoError(t, err)

		// Verify commands were parsed
		assert.Len(t, app.Commands, 5)
		assert.Equal(t, "TEMPL_EXPERIMENT=rawgo go generate ./...", app.Commands["generate"].Cmd)
		assert.Equal(t, "go build main.go", app.Commands["build"].Cmd)

		// Verify order was recorded
		require.Len(t, app.commandOrder, 5)
		assert.Equal(t, "generate", app.commandOrder[0])
		assert.Equal(t, "build", app.commandOrder[1])
		assert.Equal(t, "format", app.commandOrder[2])
		assert.Equal(t, "lint", app.commandOrder[3])
		assert.Equal(t, "test", app.commandOrder[4])

		// Marshal back to JSON
		marshaled, err := json.Marshal(&app)
		require.NoError(t, err)

		// Unmarshal again to verify order is preserved
		var app2 App
		err = json.Unmarshal(marshaled, &app2)
		require.NoError(t, err)

		// Verify order is still the same
		require.Len(t, app2.commandOrder, 5)
		assert.Equal(t, "generate", app2.commandOrder[0])
		assert.Equal(t, "build", app2.commandOrder[1])
		assert.Equal(t, "format", app2.commandOrder[2])
		assert.Equal(t, "lint", app2.commandOrder[3])
		assert.Equal(t, "test", app2.commandOrder[4])
	})

	t.Run("Default commands added in order", func(t *testing.T) {
		t.Parallel()

		app := App{
			Name: "api",
			Type: AppTypeGoLang,
			Path: "./api",
		}

		err := app.applyDefaults()
		require.NoError(t, err)

		// Verify default order was tracked
		require.Len(t, app.commandOrder, 4)
		assert.Equal(t, "format", app.commandOrder[0])
		assert.Equal(t, "lint", app.commandOrder[1])
		assert.Equal(t, "test", app.commandOrder[2])
		assert.Equal(t, "build", app.commandOrder[3])
	})

	t.Run("User commands preserved with defaults appended", func(t *testing.T) {
		t.Parallel()

		jsonData := []byte(`{
			"name": "api",
			"title": "API",
			"type": "golang",
			"path": "./api",
			"build": {},
			"infra": {"provider": "digitalocean", "type": "vm"},
			"commands": {
				"generate": {"command": "go generate ./..."},
				"build": {"command": "go build main.go"}
			}
		}`)

		var app App
		err := json.Unmarshal(jsonData, &app)
		require.NoError(t, err)

		// Apply defaults
		err = app.applyDefaults()
		require.NoError(t, err)

		// User commands should be first, defaults appended
		require.Len(t, app.commandOrder, 4)
		assert.Equal(t, "generate", app.commandOrder[0])
		assert.Equal(t, "build", app.commandOrder[1])
		// Defaults added after user commands (format and lint were added, test was not since it's default)
		assert.Contains(t, app.commandOrder, "format")
		assert.Contains(t, app.commandOrder, "lint")
		assert.Contains(t, app.commandOrder, "test")
	})
}
