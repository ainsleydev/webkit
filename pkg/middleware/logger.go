package middleware

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/util/httputil"
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
		rr := httputil.NewResponseRecorder(ctx.Response)
		ctx.Response = rr
		req := ctx.Request

		if strings.Contains(req.URL.Path, "favicon.ico") {
			return next(ctx)
		}

		if err := next(ctx); err != nil {
			return err
		}

		scheme := "http"
		if req.TLS != nil {
			scheme = "https"
		}

		level := statusLevel(rr.Status)

		msg := fmt.Sprintf("%s [%s] - %s://%s%s",
			statusLabel(rr.Status),
			strings.ToUpper(req.Method),
			scheme,
			req.Host,
			req.RequestURI,
		)

		var fields []any
		if env.IsProduction() {
			fields = []any{
				slog.String("url", req.URL.Path),
				slog.String("proto", req.Proto),
				slog.String("method", req.Method),
				slog.Int("status", rr.Status),
				slog.String("remote_addr", req.RemoteAddr),
				slog.Duration("latency", time.Now().Sub(start)),
				slog.Any(RequestIDContextKey, ctx.Get(RequestIDContextKey)),
				slog.String("user_agent", req.UserAgent()),
				slog.Any(webkit.ErrorKey, ctx.Get("error")),
				slog.Any("cache", rr.Header().Get("X-Cache")),
			}
		} else {
			msg = fmt.Sprintf("%s - %s", msg, aurora.Gray(10, fmt.Sprintf("Path: %s, Status: %d, Cache: %v, Latency: %d",
				req.URL.Path,
				rr.Status,
				rr.Header().Get("X-Cache"),
				time.Now().Sub(start).Milliseconds(),
			)))
		}

		slog.Log(ctx.Request.Context(), level, msg, fields...)

		return nil
	}
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
		//return slog.LevelInfo
		return slog.LevelError
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
