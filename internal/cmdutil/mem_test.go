package cmdutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemRunner_Run(t *testing.T) {
	t.Parallel()

	t.Run("With Stub", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		runner.AddStub("git status", Result{Output: "nothing to commit"})

		cmd := NewCommand("git", "status")
		got := runner.Run(t.Context(), cmd)

		assert.NoError(t, got.Err)
		assert.Equal(t, "git status", got.CmdLine)
		assert.Equal(t, "nothing to commit", got.Output)
	})

	t.Run("No Stub", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		cmd := NewCommand("git", "status")
		got := runner.Run(t.Context(), cmd)

		assert.Error(t, got.Err)
		assert.Contains(t, got.Err.Error(), "no stub for command")
		assert.Equal(t, "git status", got.CmdLine)
	})

	t.Run("Prefix Matching", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		runner.AddStub("git", Result{Output: "git output"})

		cmd := NewCommand("git", "status")
		got := runner.Run(t.Context(), cmd)

		assert.NoError(t, got.Err)
		assert.Equal(t, "git output", got.Output)
		assert.Equal(t, "git status", got.CmdLine)
	})

	t.Run("With Error", func(t *testing.T) {
		t.Parallel()

		runner := NewMemRunner()
		runner.AddStub("fail", Result{
			Output: "error output",
			Err:    errors.New("command failed"),
		})

		cmd := NewCommand("fail")
		got := runner.Run(t.Context(), cmd)

		assert.Error(t, got.Err)
		assert.Contains(t, got.Err.Error(), "command failed")
		assert.Equal(t, "error output", got.Output)
		assert.Equal(t, "fail", got.CmdLine)
	})
}

func TestMemRunner_Calls(t *testing.T) {
	t.Parallel()

	runner := NewMemRunner()
	runner.AddStub("echo", Result{Output: "test"})

	cmd1 := NewCommand("echo", "hello")
	cmd2 := NewCommand("echo", "world")

	runner.Run(t.Context(), cmd1)
	runner.Run(t.Context(), cmd2)

	calls := runner.Calls()

	assert.Len(t, calls, 2)
	assert.Equal(t, "echo hello", calls[0].String())
	assert.Equal(t, "echo world", calls[1].String())
}

func TestMemRunner_Reset(t *testing.T) {
	t.Parallel()

	runner := NewMemRunner()
	runner.AddStub("test", Result{Output: "output"})

	cmd := NewCommand("test", "arg")
	runner.Run(t.Context(), cmd)

	assert.Len(t, runner.Calls(), 1)

	runner.Reset()

	assert.Len(t, runner.Calls(), 0)

	// Stubs should be cleared too
	result := runner.Run(t.Context(), cmd)
	assert.Error(t, result.Err)
}
