package cmdutil

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// MemRunner implements Runner by stubbing calls when
// a command is run.
type MemRunner struct {
	mu    sync.Mutex
	calls []Command
	stubs map[string]Result
}

// NewMemRunner creates a new runner that stubs calls.
func NewMemRunner() *MemRunner {
	return &MemRunner{
		stubs: make(map[string]Result),
	}
}

func (r *MemRunner) Run(_ context.Context, cmd Command) Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.calls = append(r.calls, cmd)

	cmdLine := cmd.String()
	for prefix, stub := range r.stubs {
		if strings.HasPrefix(cmdLine, prefix) {
			// Always set the actual command line
			stub.CmdLine = cmdLine
			return stub
		}
	}

	return Result{
		CmdLine: cmdLine,
		Output:  "",
		Err:     fmt.Errorf("no stub for command: %s", cmdLine),
	}
}

// AddStub registers a command prefix and a fake result.
func (r *MemRunner) AddStub(prefix string, result Result) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stubs[prefix] = result
}

// Calls returns all commands that were executed.
func (r *MemRunner) Calls() []Command {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Command(nil), r.calls...)
}

// Reset clears recorded calls.
func (r *MemRunner) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = nil
	r.stubs = make(map[string]Result)
}
