package manifest

import (
	"runtime"
	"strings"
	"unicode"
)

// CallerInfo holds metadata about the detected caller.
type CallerInfo struct {
	Package  string // e.g. "env"
	Function string // e.g. "Scaffold"
	File     string // e.g. "/path/to/env.go"
	Line     int    // e.g. 42
}

// Caller returns information about the first exported (public)
// function in the call stack outside of runtime/scaffold packages.
func Caller() CallerInfo {
	const maxDepth = 15
	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		name := frame.Function // e.g., github.com/ainsleydev/webkit/internal/env.Scaffold
		if strings.Contains(name, "runtime.") || strings.Contains(name, "/scaffold.") {
			continue
		}

		parts := strings.Split(name, "/")
		last := parts[len(parts)-1]         // e.g. "env.Scaffold"
		fnParts := strings.Split(last, ".") // ["env", "Scaffold"]
		if len(fnParts) < 2 {
			continue
		}
		pkg := fnParts[len(fnParts)-2]
		fn := fnParts[len(fnParts)-1]

		// Only return if the function is exported (starts with uppercase).
		if len(fn) > 0 && unicode.IsUpper(rune(fn[0])) {
			return CallerInfo{
				Package:  pkg,
				Function: fn,
				File:     frame.File,
				Line:     frame.Line,
			}
		}
	}

	return CallerInfo{
		Package:  "unknown",
		Function: "unknown",
	}
}
