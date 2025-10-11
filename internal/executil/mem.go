package executil

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
	errs  map[string]error
}

// NewMemRunner creates a new runner that stubs calls.
func NewMemRunner() *MemRunner {
	return &MemRunner{
		stubs: make(map[string]Result),
		errs:  make(map[string]error),
	}
}

func (r *MemRunner) Run(_ context.Context, cmd Command) (Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.calls = append(r.calls, cmd)

	cmdLine := cmd.String()
	for prefix, stub := range r.stubs {
		if strings.HasPrefix(cmdLine, prefix) {
			stub.CmdLine = cmdLine
			return stub, r.errs[prefix] // Returns error if set
		}
	}

	return Result{
		CmdLine: cmdLine,
		Output:  "",
	}, fmt.Errorf("no stub for command: %s", cmdLine)
}

// AddStub registers a command prefix and a fake result.
func (r *MemRunner) AddStub(prefix string, result Result, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stubs[prefix] = result
	if err != nil {
		r.errs[prefix] = err
	} else {
		delete(r.errs, prefix)
	}
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
