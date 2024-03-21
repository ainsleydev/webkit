package log

import (
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

func resolveLogHandler(prefix string) slog.Handler {
	enviro := env.AppEnvironment()
	if enviro == env.Staging || enviro == env.Production {
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil, // We don't need to replace anything at the moment, but sensitive info could be masked.
		})
	}
	return NewLocalHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}, prefix)
}
