package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white. Logger prints a
// request ID if one is provided.
//
// IMPORTANT NOTE: Logger should go before any other middleware that may change
// the response, such as middleware.Recover. Example:
//
//	app := webkit.New()
//	app.Plug(middleware.Logger)        // <--<< Logger should come before Recover
//	app.Plug(middleware.Recover)
//	r.Get("/", handler)
func Logger(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		start := time.Now()
		rw := &responseWrapper{ResponseWriter: ctx.Response}
		req := ctx.Request

		if strings.Contains(req.URL.Path, "favicon.ico") {
			return next(ctx)
		}

		next.ServeHTTP(rw, ctx.Request)

		scheme := "http"
		if req.TLS != nil {
			scheme = "https"
		}

		level := statusLevel(rw.status)

		msg := fmt.Sprintf("%s [%s] - %s://%s%s %s",
			statusLabel(rw.status),
			strings.ToUpper(req.Method),
			scheme,
			req.Host,
			req.RequestURI,
			req.Proto,
		)

		fields := []any{
			slog.String("url", req.URL.Path),
			slog.String("method", req.Method),
			slog.Int("status", rw.status),
			slog.String("remote_addr", req.RemoteAddr),
			slog.Duration("latency", time.Now().Sub(start)),
			slog.Any(RequestIDContextKey, ctx.Get(RequestIDContextKey)),
			slog.String("user_agent", req.UserAgent()),
			slog.Any(webkit.ErrorKey, ctx.Get("error")),
		}

		if !env.IsProduction() {
			fields = fields[2:]
		}

		slog.Log(ctx.Context(), level, msg, fields...)

		return nil
	}
}

// responseWrapper is a struct that wraps a http.ResponseWriter to intercept the status code
type responseWrapper struct {
	http.ResponseWriter
	status int
}

// WriteHeader intercepts and stores the status code before writing it to the client
func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// statusLevel returns a slog.Level based on the HTTP status code.
func statusLevel(status int) slog.Level {
	switch {
	case status <= 0:
		return slog.LevelWarn
	case status < 400: // for codes in 100s, 200s, 300s
		return slog.LevelInfo
	case status >= 400 && status < 500:
		// switching to info level to be less noisy
		return slog.LevelInfo
	case status >= 500:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// statusLabel returns a human readable status code label.
func statusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return "OK"
	case status >= 300 && status < 400:
		return "Redirect"
	case status >= 400 && status < 500:
		return "Client Error"
	case status >= 500:
		return "Server Error"
	default:
		return "Unknown"
	}
}
