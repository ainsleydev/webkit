package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandSpecUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   []byte
		want    CommandSpec
		wantErr bool
	}{
		"Bool True": {
			input: []byte(`true`),
			want:  CommandSpec{Disabled: false},
		},
		"Bool False": {
			input: []byte(`false`),
			want:  CommandSpec{Disabled: true},
		},
		"String Override": {
			input: []byte(`"custom command"`),
			want:  CommandSpec{Cmd: "custom command"},
		},
		"Full Object": {
			input: []byte(`{"command":"run tests","skip_ci":true,"timeout":"5m"}`),
			want:  CommandSpec{Cmd: "run tests", SkipCI: true, Timeout: "5m"},
		},
		"Invalid JSON": {
			input:   []byte(`{"invalid":`),
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got CommandSpec
			err := got.UnmarshalJSON(test.input)

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestAppGetCommand(t *testing.T) {
	t.Parallel()

	var commands = []Command{
		CommandLint,
		CommandTest,
		CommandFormat,
	}

	t.Run("Defaults", func(t *testing.T) {
		t.Parallel()

		for appType, expectedCmds := range defaultCommands {
			t.Run("App Type "+string(appType), func(t *testing.T) {
				t.Parallel()

				app := App{Type: appType}

				for _, command := range commands {
					expected := expectedCmds[command]

					t.Run(string(command), func(t *testing.T) {
						t.Parallel()
						got, skip := app.GetCommand(command)
						assert.Equal(t, expected, got)
						assert.False(t, skip, "Default command should not be skippable")
						assert.NotEmpty(t, got, "Default command should not be empty")
					})
				}
			})
		}
	})

	tt := map[string]struct {
		app   App
		input Command
		want  string
		skip  bool
	}{
		"Default Command Exists": {
			app: App{
				Type: AppTypeGoLang,
			},
			input: CommandLint,
			want:  "golangci-lint run",
			skip:  false,
		},
		"User Override Command": {
			app: App{
				Type: AppTypeSvelteKit,
				Commands: map[Command]CommandSpec{
					CommandTest: {Cmd: "custom test", SkipCI: true},
				},
			},
			input: CommandTest,
			want:  "custom test",
			skip:  true,
		},
		"User Disabled Command": {
			app: App{
				Type: AppTypePayload,
				Commands: map[Command]CommandSpec{
					CommandFormat: {Disabled: true},
				},
			},
			input: CommandFormat,
			want:  "",
			skip:  true,
		},
		"User Override Empty Cmd Fallbacks To Default": {
			app: App{
				Type: AppTypeGoLang,
				Commands: map[Command]CommandSpec{
					CommandTest: {Cmd: "", SkipCI: true},
				},
			},
			input: CommandTest,
			want:  "go test ./...",
			skip:  true,
		},
		"No Defaults For Type": {
			app: App{
				Type: "unknown-type",
			},
			input: CommandLint,
			want:  "",
			skip:  true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, skip := test.app.GetCommand(test.input)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.skip, skip)
		})
	}
}
