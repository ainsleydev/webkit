package cmdutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ExecRunner implements Runner by using cmd.Execute to
// run a Command.
type ExecRunner struct{}

// DefaultRunner returns a real command runner.
func DefaultRunner() Runner {
	return &ExecRunner{}
}

func (e ExecRunner) Run(ctx context.Context, cmd Command) Result {
	result := Result{
		CmdLine: cmd.String(),
	}

	// Prepare command
	c := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	c.Dir = cmd.Dir
	c.Stdin = cmd.Stdin

	// Merge env
	env := os.Environ()
	for k, v := range cmd.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	c.Env = env

	var stdoutBuf, stderrBuf bytes.Buffer

	// Set output
	if cmd.Stdout != nil {
		c.Stdout = io.MultiWriter(cmd.Stdout, &stdoutBuf)
	} else {
		c.Stdout = &stdoutBuf
	}

	if cmd.Stderr != nil {
		c.Stderr = io.MultiWriter(cmd.Stderr, &stderrBuf)
	} else {
		c.Stderr = &stderrBuf
	}

	// Run
	err := c.Run()
	result.Err = err
	result.Output = stdoutBuf.String() + stderrBuf.String()
	return result
}
