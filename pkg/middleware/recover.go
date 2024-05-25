package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Recover is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. The error is also reported to Sentry.
//
// This middleware should be plugged in first to ensure that it catches any
// panics that occur in the request-response cycle.
func Recover(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().RecoverWithContext(ctx.Request.Context(), err)
				sentry.Flush(time.Second * 5)

				slog.ErrorContext(ctx.Request.Context(), "Panic recovered", slog.Any("error", err))

				if ctx.Request.Header.Get("Connection") != "Upgrade" {
					ctx.Response.WriteHeader(http.StatusInternalServerError)
				}

				panic(err)
			}
		}()
		return next(ctx)
	}
}
