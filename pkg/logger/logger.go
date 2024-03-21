package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/ainsleydev/webkit/pkg/env"
)

// Bootstrap creates a new complaint logger and sets the default.
func Bootstrap(prefix string) {
	slog.SetDefault(slog.New(resolveLogHandler(prefix)))
}

// Debug calls [Logger.Debug] on the default logger.
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// DebugContext calls [Logger.DebugContext] on the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// Info calls [Logger.Info] on the default logger.
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// InfoContext calls [Logger.InfoContext] on the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// Warn calls [Logger.Warn] on the default logger.
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// WarnContext calls [Logger.WarnContext] on the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// Error calls [Logger.Error] on the default logger.
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// ErrorContext calls [Logger.ErrorContext] on the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// Log calls [Logger.Log] on the default logger.
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}

// LogAttrs calls [Logger.LogAttrs] on the default logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, level, msg, attrs...)
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
