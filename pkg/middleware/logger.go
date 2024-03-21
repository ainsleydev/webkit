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
// Alternatively, look at https://github.com/goware/httplog for a more in-depth
// http logger with structured logging support.
//
// IMPORTANT NOTE: Logger should go before any other middleware that may change
// the response, such as middleware.Recoverer. Example:
//
//	r := chi.NewRouter()
//	r.Use(middleware.Logger)        // <--<< Logger should come before Recoverer
//	r.Use(middleware.Recoverer)
//	r.Get("/", handler)
func Logger(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		start := time.Now()
		rw := &responseWrapper{ResponseWriter: ctx.Response}
		req := ctx.Request

		next.ServeHTTP(rw, ctx.Request)

		scheme := "http"
		if req.TLS != nil {
			scheme = "https"
		}

		// TODO: Think we need to add a coloured pretty print of the status code.
		msg := fmt.Sprintf("[%s] - %s://%s%s %s", strings.ToUpper(req.Method), scheme, req.Host, req.RequestURI, req.Proto)

		fields := []any{
			slog.String("url", req.URL.Path),
			slog.String("method", req.Method),
			slog.Int("status", rw.status),
			slog.String("remote_addr", req.RemoteAddr),
			slog.Duration("latency", time.Now().Sub(start)),
			slog.Any("request_id", ctx.Get(RequestIDContextKey)),
		}

		if env.IsProduction() {
			slog.DebugContext(ctx.Context(), msg, fields...)
			return nil
		}

		slog.Info(msg, fields[2:]...)

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
