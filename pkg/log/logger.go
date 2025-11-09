package log

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/ainsleydev/webkit/pkg/env"
)

// DefaultLogger is the global logger instance used by WebKit for HTTP requests
// and application-wide logging. It can be configured with Bootstrap.
var DefaultLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// Bootstrap initializes the default logger with environment-aware formatting.
// The prefix appears in development logs to identify the application.
// Call this once during application startup.
func Bootstrap(prefix string) {
	DefaultLogger = slog.New(resolveLogHandler(prefix))
	slog.SetDefault(DefaultLogger)
}

// NewLogger creates a new logger instance with environment-aware formatting.
// Use this to create isolated loggers with different prefixes for subsystems.
func NewLogger(prefix string) *slog.Logger {
	return slog.New(resolveLogHandler(prefix))
}

// NewNoOpLogger creates a logger that discards all output.
// Useful for testing or when logging should be completely suppressed.
func NewNoOpLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// Discard configures the standard library logger to discard all output.
// This only affects log.Print calls, not slog or DefaultLogger.
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
