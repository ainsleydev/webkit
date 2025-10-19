package executil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemRunner_Run(t *testing.T) {
	t.Parallel()

	t.Run("With Stub", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		runner.AddStub("git status", Result{Output: "nothing to commit"}, nil)
		cmd := NewCommand("git", "status")

		got, err := runner.Run(t.Context(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, "git status", got.CmdLine)
		assert.Equal(t, "nothing to commit", got.Output)
	})

	t.Run("No Stub", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		cmd := NewCommand("git", "status")

		got, err := runner.Run(t.Context(), cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no stub for command")
		assert.Equal(t, "git status", got.CmdLine)
	})

	t.Run("Prefix Matching", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		runner.AddStub("git", Result{Output: "git output"}, nil)
		cmd := NewCommand("git", "status")

		got, err := runner.Run(t.Context(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, "git output", got.Output)
		assert.Equal(t, "git status", got.CmdLine)
	})

	t.Run("With Error", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		runner.AddStub("fail", Result{
			Output: "error output",
		}, fmt.Errorf("command failed"))

		cmd := NewCommand("fail")
		got, err := runner.Run(t.Context(), cmd)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "command failed")
		assert.Equal(t, "error output", got.Output)
		assert.Equal(t, "fail", got.CmdLine)
	})
}

func TestMemRunner_Calls(t *testing.T) {
	t.Parallel()

	runner := NewMemRunner()
	runner.AddStub("echo", Result{Output: "test"}, nil)

	cmd1 := NewCommand("echo", "hello")
	cmd2 := NewCommand("echo", "world")

	_, err := runner.Run(t.Context(), cmd1)
	require.NoError(t, err)
	_, err = runner.Run(t.Context(), cmd2)
	require.NoError(t, err)

	calls := runner.Calls()

	assert.Len(t, calls, 2)
	assert.Equal(t, "echo hello", calls[0].String())
	assert.Equal(t, "echo world", calls[1].String())
}

func TestMemRunner_Reset(t *testing.T) {
	t.Parallel()

	runner := NewMemRunner()
	runner.AddStub("test", Result{Output: "output"}, nil)

	cmd := NewCommand("test", "arg")
	_, err := runner.Run(t.Context(), cmd)
	assert.NoError(t, err)

	assert.Len(t, runner.Calls(), 1)

	runner.Reset()

	assert.Len(t, runner.Calls(), 0)

	// Stubs should be cleared too
	_, err = runner.Run(t.Context(), cmd)
	assert.Error(t, err)
}
