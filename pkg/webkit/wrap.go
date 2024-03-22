package webkit

import "net/http"

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
