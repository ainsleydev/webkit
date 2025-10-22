package executil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_String(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Command
		want  string
	}{
		"No Args": {
			input: NewCommand("git"),
			want:  "git",
		},
		"Single Arg": {
			input: NewCommand("echo", "hello"),
			want:  "echo hello",
		},
		"Multiple Args": {
			input: NewCommand("build", "-o", "bin/app"),
			want:  "build -o bin/app",
		},
		"Empty Key": {
			input: NewCommand(""),
			want:  "",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func Test_Exists(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  bool
	}{
		"Existing Command": {
			input: "echo",
			want:  true,
		},
		"Nonexistent Command": {
			input: "nonexistentcommand123",
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := Exists(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
