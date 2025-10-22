package executil

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecRunner_Run(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	tt := map[string]struct {
		input       Command
		wantOutput  []string
		wantCmdLine string
	}{
		"Simple": {
			input:       NewCommand("echo", "hello", "world"),
			wantOutput:  []string{"hello world"},
			wantCmdLine: "echo hello world",
		},
		"With Stdout": {
			input: func() Command {
				var buf bytes.Buffer
				c := NewCommand("echo", "test")
				c.Stdout = &buf
				return c
			}(),
			wantOutput:  []string{"test"},
			wantCmdLine: "echo test",
		},
		"With Stderr": {
			input: func() Command {
				var buf bytes.Buffer
				c := NewCommand("sh", "-c", "echo error >&2")
				c.Stderr = &buf
				return c
			}(),
			wantOutput:  []string{"error"},
			wantCmdLine: "sh -c echo error >&2",
		},
		"With Env": {
			input: func() Command {
				c := NewCommand("sh", "-c", "echo $TEST_VAR")
				c.Env = map[string]string{"TEST_VAR": "hello"}
				return c
			}(),
			wantOutput:  []string{"hello"},
			wantCmdLine: "sh -c echo $TEST_VAR",
		},
		"With Dir": {
			input: func() Command {
				c := NewCommand("pwd")
				c.Dir = tmpDir
				return c
			}(),
			wantOutput:  []string{tmpDir},
			wantCmdLine: "pwd",
		},
		"With Stdin": {
			input: func() Command {
				c := NewCommand("cat")
				c.Stdin = strings.NewReader("test input\n")
				return c
			}(),
			wantOutput:  []string{"test input"},
			wantCmdLine: "cat",
		},
		"With Stdout & Stderr": {
			input: func() Command {
				var stdoutBuf, stderrBuf bytes.Buffer
				c := NewCommand("sh", "-c", "echo stdout; echo stderr >&2")
				c.Stdout = &stdoutBuf
				c.Stderr = &stderrBuf
				return c
			}(),
			wantOutput:  []string{"stdout", "stderr"},
			wantCmdLine: "sh -c echo stdout; echo stderr >&2",
		},
		"Multiple Env Vars": {
			input: func() Command {
				c := NewCommand("sh", "-c", "echo $VAR1 $VAR2")
				c.Env = map[string]string{"VAR1": "hello", "VAR2": "world"}
				return c
			}(),
			wantOutput:  []string{"hello", "world"},
			wantCmdLine: "sh -c echo $VAR1 $VAR2",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := DefaultRunner()

			got, err := runner.Run(t.Context(), test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.wantCmdLine, got.CmdLine)
			for _, w := range test.wantOutput {
				assert.Contains(t, got.Output, w)
			}
		})
	}

	t.Run("Bad Exit Code", func(t *testing.T) {
		t.Parallel()

		runner := DefaultRunner()
		cmd := NewCommand("sh", "-c", "exit 42")

		_, err := runner.Run(t.Context(), cmd)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "exit status 42")
	})

	t.Run("Command Not Found", func(t *testing.T) {
		t.Parallel()

		runner := DefaultRunner()
		cmd := NewCommand("this-command-does-not-exist")

		_, err := runner.Run(t.Context(), cmd)
		assert.Error(t, err)
	})
}
