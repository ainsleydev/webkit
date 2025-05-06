package log

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/ainsleydev/webkit/pkg/env"
)

// DefaultLogger is the logger that WebKit uses to log HTTP Requests and
// info messages throughout the application.
var DefaultLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// Bootstrap creates a new complaint logger and sets the default.
func Bootstrap(prefix string) {
	DefaultLogger = slog.New(resolveLogHandler(prefix))
	slog.SetDefault(DefaultLogger)
}

// NewLogger creates a new WebKit compliant logger with the given prefix.
func NewLogger(prefix string) *slog.Logger {
	return slog.New(resolveLogHandler(prefix))
}

// NewNoOpLogger creates a logger that does nothing.
func NewNoOpLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// Discard sets the default logger to discard all logs.
// It's an alias for log.SetOutput(io.Discard).
func Discard() {
	log.SetOutput(io.Discard)
}

func resolveLogHandler(prefix string) slog.Handler {
	if !env.IsDevelopment() {
		return jsonHandler{
			Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				AddSource:   false,
				Level:       slog.LevelInfo,
				ReplaceAttr: nil, // We don't need to replace anything at the moment, but sensitive info could be masked.
			}),
		}
	}
	return NewLocalHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}, prefix)
}
