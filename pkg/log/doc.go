// Package log provides structured logging for WebKit applications.
//
// It wraps Go's standard log/slog package with environment-aware handlers
// that produce human-readable output in development and structured JSON
// in production. The package automatically includes request IDs from context
// when available.
//
// Basic usage:
//
//	log.Bootstrap("myapp")
//	slog.Info("server started", "port", 8080)
//
// The Bootstrap function configures the default logger based on environment.
// In development, logs are formatted with colors and timestamps for readability.
// In production, logs are output as structured JSON.
package log
