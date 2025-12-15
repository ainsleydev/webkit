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
		"Full Object with WorkingDirectory": {
			input: []byte(`{"command":"run tests","skip_ci":true,"timeout":"5m","working_directory":"./subdir"}`),
			want:  CommandSpec{Cmd: "run tests", SkipCI: true, Timeout: "5m", WorkingDirectory: "./subdir"},
		},
		"Object with only WorkingDirectory": {
			input: []byte(`{"command":"pnpm build","working_directory":"./packages/core"}`),
			want:  CommandSpec{Cmd: "pnpm build", WorkingDirectory: "./packages/core"},
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

func TestCommand_String(t *testing.T) {
	t.Parallel()

	got := CommandLint.String()
	assert.Equal(t, "lint", got)
	assert.IsType(t, "", got)
}
