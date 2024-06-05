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

func WrapKitHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(NewContext(w, r))
	}
}

// WrapMiddleware wraps `func(http.Handler) http.Handler` into `webkit.Plug“
func WrapMiddleware(m func(http.Handler) http.Handler) Plug {
	return func(next Handler) Handler {
		return func(c *Context) (err error) {
			m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err = next(c)
			})).ServeHTTP(c.Response, c.Request)
			return
		}
	}
}

func WrapMiddelewareHandler(next Handler, m func(http.Handler) http.Handler) Handler {
	return func(c *Context) (err error) {
		m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Response = w
			c.Request = r
			err = next(c)
		})).ServeHTTP(c.Response, c.Request)
		return
	}
}
