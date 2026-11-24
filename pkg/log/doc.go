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
// In production, logs are output as structured JSON, one object per line,
// with every attribute written once at the top level. ERROR and FATAL
// records go to stderr and everything else to stdout, because platforms
// such as Vercel label log lines by stream. Use NewJSONHandler to change
// the writers or the stderr threshold.
//
// To report errors to Sentry, wrap the handler with NewSentryHandler:
//
//	h := log.NewSentryHandler(log.NewHandler("myapp"), sentryClient, nil)
//	slog.SetDefault(slog.New(h))
//	slog.Error("payment failed", "error", err) // Also sent to Sentry.
//
// NewSentryHandler replaces the logrus-based SentryHook, which is
// deprecated and will be removed in a future release.
package log
