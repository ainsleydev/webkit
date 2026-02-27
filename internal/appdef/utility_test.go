package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef/types"
)

func TestUtility_HasCI(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		ci   *UtilityCI
		want bool
	}{
		"Nil CI":     {ci: nil, want: false},
		"Non-nil CI": {ci: &UtilityCI{Trigger: "pull_request"}, want: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			u := Utility{CI: test.ci}
			got := u.HasCI()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestUtility_OrderedCommands(t *testing.T) {
	t.Parallel()

	t.Run("Nil Commands Returns Nil", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "e2e",
			Path:     "e2e",
			Language: "js",
		}

		commands := u.OrderedCommands()
		assert.Nil(t, commands)
	})

	t.Run("Empty Commands", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "e2e",
			Path:     "e2e",
			Language: "js",
			Toolset:  Toolset{Commands: types.NewOrderedMap[Command, CommandSpec]()},
		}

		commands := u.OrderedCommands()
		assert.Len(t, commands, 0)
	})

	t.Run("Populated Commands Preserve Order", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "e2e",
			Path:     "e2e",
			Language: "js",
			Toolset:  Toolset{Commands: types.NewOrderedMap[Command, CommandSpec]()},
		}

		u.Commands.Set("test", CommandSpec{Cmd: "pnpm playwright test"})
		u.Commands.Set("report", CommandSpec{Cmd: "pnpm playwright show-report"})

		commands := u.OrderedCommands()
		require.Len(t, commands, 2)
		assert.Equal(t, "test", commands[0].Name)
		assert.Equal(t, "pnpm playwright test", commands[0].Cmd)
		assert.Equal(t, "report", commands[1].Name)
		assert.Equal(t, "pnpm playwright show-report", commands[1].Cmd)
	})
}

func TestUtility_ShouldUseNPM(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		language string
		want     bool
	}{
		"JS": {language: "js", want: true},
		"Go": {language: "go", want: false},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			u := Utility{Language: test.language}
			got := u.ShouldUseNPM()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestUtility_InstallCommands(t *testing.T) {
	t.Parallel()

	t.Run("Go tools", func(t *testing.T) {
		t.Parallel()

		u := Utility{
			Language: "go",
			Toolset: Toolset{Tools: map[string]Tool{
				"golangci-lint": {Type: "go", Name: "github.com/golangci/golangci-lint/cmd/golangci-lint", Version: "latest"},
			}},
		}

		got := u.InstallCommands()
		assert.Len(t, got, 1)
		assert.Contains(t, got, "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	})

	t.Run("Pnpm tools", func(t *testing.T) {
		t.Parallel()

		u := Utility{
			Language: "js",
			Toolset: Toolset{Tools: map[string]Tool{
				"playwright": {Type: "pnpm", Name: "@playwright/test", Version: "latest"},
			}},
		}

		got := u.InstallCommands()
		assert.Len(t, got, 1)
		assert.Contains(t, got, "pnpm add -g @playwright/test@latest")
	})

	t.Run("Script tools with install override", func(t *testing.T) {
		t.Parallel()

		u := Utility{
			Language: "js",
			Toolset: Toolset{Tools: map[string]Tool{
				"custom": {
					Type:    "script",
					Install: "curl -sSL https://example.com/install.sh | sh",
				},
			}},
		}

		got := u.InstallCommands()
		assert.Contains(t, got, "curl -sSL https://example.com/install.sh | sh")
	})

	t.Run("Deterministic ordering", func(t *testing.T) {
		t.Parallel()

		u := Utility{
			Language: "go",
			Toolset: Toolset{Tools: map[string]Tool{
				"zebra": {Type: "go", Name: "github.com/z/zebra", Version: "v1.0.0"},
				"alpha": {Type: "go", Name: "github.com/a/alpha", Version: "v1.0.0"},
			}},
		}

		got := u.InstallCommands()
		want := []string{
			"go install github.com/a/alpha@v1.0.0",
			"go install github.com/z/zebra@v1.0.0",
		}
		assert.Equal(t, want, got)
	})

	t.Run("Empty tools", func(t *testing.T) {
		t.Parallel()

		u := Utility{Language: "js", Toolset: Toolset{Tools: map[string]Tool{}}}
		got := u.InstallCommands()
		assert.Nil(t, got)
	})
}

func TestUtility_ApplyDefaults(t *testing.T) {
	t.Parallel()

	t.Run("Initialises nil Commands and Tools", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "e2e",
			Path:     "e2e",
			Language: "js",
		}

		err := u.applyDefaults()
		assert.NoError(t, err)

		require.NotNil(t, u.Commands)
		require.NotNil(t, u.Tools)
		assert.Len(t, u.Tools, 0)
	})

	t.Run("Cleans path", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "e2e",
			Path:     "./tests/../e2e",
			Language: "js",
		}

		err := u.applyDefaults()
		assert.NoError(t, err)

		assert.Equal(t, "e2e", u.Path)
	})

	t.Run("Defaults CI RunsOn", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "e2e",
			Path:     "e2e",
			Language: "js",
			CI:       &UtilityCI{Trigger: "pull_request"},
		}

		err := u.applyDefaults()
		assert.NoError(t, err)

		assert.Equal(t, "ubuntu-latest", u.CI.RunsOn)
	})

	t.Run("Preserves explicit CI RunsOn", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "e2e",
			Path:     "e2e",
			Language: "js",
			CI:       &UtilityCI{Trigger: "pull_request", RunsOn: "ubuntu-22.04"},
		}

		err := u.applyDefaults()
		assert.NoError(t, err)

		assert.Equal(t, "ubuntu-22.04", u.CI.RunsOn)
	})

	t.Run("No CI defaults when CI is nil", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "constants",
			Path:     "packages/constants",
			Language: "js",
		}

		err := u.applyDefaults()
		assert.NoError(t, err)

		assert.Nil(t, u.CI)
	})

	t.Run("Cron trigger with schedule", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "nightly-e2e",
			Path:     "e2e",
			Language: "js",
			CI:       &UtilityCI{Trigger: "cron", Schedule: "0 0 * * *"},
		}

		err := u.applyDefaults()
		assert.NoError(t, err)

		assert.Equal(t, "ubuntu-latest", u.CI.RunsOn)
		assert.Equal(t, "0 0 * * *", u.CI.Schedule)
	})

	t.Run("Cron trigger without schedule returns error", func(t *testing.T) {
		t.Parallel()

		u := &Utility{
			Name:     "nightly-e2e",
			Path:     "e2e",
			Language: "js",
			CI:       &UtilityCI{Trigger: "cron"},
		}

		err := u.applyDefaults()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cron trigger but no schedule")
	})

	t.Run("Preserves existing Commands and Tools", func(t *testing.T) {
		t.Parallel()

		cmds := types.NewOrderedMap[Command, CommandSpec]()
		cmds.Set("test", CommandSpec{Cmd: "pnpm test"})

		u := &Utility{
			Name:     "e2e",
			Path:     "e2e",
			Language: "js",
			Toolset: Toolset{
				Tools: map[string]Tool{
					"playwright": {Type: "pnpm", Name: "@playwright/test", Version: "latest"},
				},
				Commands: cmds,
			},
		}

		err := u.applyDefaults()
		assert.NoError(t, err)

		assert.Len(t, u.Tools, 1)
		spec, exists := u.Commands.Get("test")
		require.True(t, exists)
		assert.Equal(t, "pnpm test", spec.Cmd)
	})
}
