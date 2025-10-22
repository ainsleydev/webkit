package executil

import (
	"context"
	"io"
	"os/exec"
	"strings"
)

type (
	// Runner is the interface for running external commands
	// using exec.Command.
	Runner interface {
		Run(ctx context.Context, cmd Command) (Result, error)
	}
	// Command defines all options for running a command.
	Command struct {
		Name   string
		Args   []string
		Dir    string
		Env    map[string]string
		Stdin  io.Reader
		Stdout io.Writer
		Stderr io.Writer
	}
	// Result captures the outcome of running a command.
	Result struct {
		CmdLine string
		Output  string
	}
)

// NewCommand creates a new Command with a name and args.
func NewCommand(name string, args ...string) Command {
	return Command{Name: name, Args: args}
}

// String implements fmt.Stringer on the command to print
// the name of the command and it's arguments.
func (c Command) String() string {
	if len(c.Args) == 0 {
		return c.Name
	}
	return c.Name + " " + strings.Join(c.Args, " ")
}

// Exists checks to see if a command line application
// exists named by the PATH environment variable.
func Exists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
