package webkit

import (
	"net/http"
)

// WrapHandlerFunc is a helper function for wrapping http.HandlerFunc
// and returns a WebKit middleware.
func WrapHandlerFunc(f http.HandlerFunc) Handler {
	return func(ctx *Context) error {
		f(ctx.Response, ctx.Request)
		return nil
	}
}

// WrapHandler is a helper function for wrapping http.Handler
// and returns a WebKit middleware.
func WrapHandler(h http.Handler) Handler {
	return func(ctx *Context) error {
		h.ServeHTTP(ctx.Response, ctx.Request)
		return nil
	}
}

// WrapMiddleware wraps `func(http.Handler) http.Handler` into `webkit.Plug“
func WrapMiddleware(next Handler, mw func(http.Handler) http.Handler) Handler {
	return func(ctx *Context) (err error) {
		mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Response = w
			ctx.Request = r
			err = next(ctx) // Call the next handler with updated context
		})).ServeHTTP(ctx.Response, ctx.Request)
		return
	}
}
