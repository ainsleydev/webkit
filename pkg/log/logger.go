package log

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/ainsleydev/webkit/pkg/env"
)

// LevelFatal is a level above slog.LevelError for failures the
// application can't recover from. Log at it with slog.Log, e.g.
//
//	slog.Log(ctx, log.LevelFatal, "database unreachable", "error", err)
//
// Handlers in this package write it as "FATAL".
const LevelFatal = slog.Level(12)

const levelFatalName = "FATAL"

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

// NewHandler returns the environment-aware handler that Bootstrap and
// NewLogger use. Use it to wrap the handler, for example with
// NewSentryHandler:
//
//	h := log.NewSentryHandler(log.NewHandler("myapp"), client, nil)
//	slog.SetDefault(slog.New(h))
func NewHandler(prefix string) slog.Handler {
	return resolveLogHandler(prefix)
}

func resolveLogHandler(prefix string) slog.Handler {
	if !env.IsDevelopment() {
		// ERROR and above go to stderr so that platforms which label
		// lines by stream (e.g. Vercel) show them at the right level.
		return NewJSONHandler(os.Stdout, os.Stderr, &JSONHandlerOptions{
			HandlerOptions: slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
			},
			StderrLevel: slog.LevelError,
		})
	}
	return NewLocalHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}, prefix)
}
