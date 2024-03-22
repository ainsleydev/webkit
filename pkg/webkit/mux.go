package webkit

import (
	"log/slog"
	"net/http"
	"strings"
)

type (
	// Kit is the top-level framework instance for handling HTTP routes.
	//
	// Goroutine safety: Do not mutate WebKit instance fields after server has started. Accessing these
	// fields from handlers/middlewares and changing field values at the same time leads to data-races.
	// Adding new routes after the server has been started is also not safe!
	Kit struct {
		ErrorHandler ErrorHandler
		mux          *http.ServeMux
		routes       []Route
		plugs        []Plug
	}
	// Route contains a handler and information for matching against requests.
	Route struct {
		Path   string
		Method string
		Name   string
	}
	// Handler is a function that handles HTTP requests.
	Handler func(c *Context) error
	// Plug defines a function to process middleware.
	Plug func(handler Handler) Handler
	// ErrorHandler is a centralized HTTP error handler.
	ErrorHandler func(*Context, error) error
)

// New creates a new WebKit instance.
func New() *Kit {
	return &Kit{
		ErrorHandler: DefaultErrorHandler,
		mux:          http.NewServeMux(),
		plugs:        []Plug{},
	}
}

// ServeHTTP implements the http.Handler interface.
func (a *Kit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// ServeHTTP wraps the Handler function to a Handler so that
// it satisfies the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(NewContext(w, r))
	if err != nil {
		slog.Error("Handling HTTP route: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Plug adds a middleware function to the chain. These are called after
// the funcs that are passed directly to the route-level handlers.
//
// For example: app.Plug(middleware.Logger)
func (a *Kit) Plug(plugs ...Plug) {
	a.plugs = append(a.plugs, plugs...)
}

// Start starts the HTTP server.
func (a *Kit) Start(address string) error {
	slog.Info("App listening on address: " + address)
	return http.ListenAndServe(address, a.mux)
}

// ErrorKey is the key used to store the error in the context.
const ErrorKey = "error"

// DefaultErrorHandler is the default error handler that is called when a route
// handler returns an error.
var DefaultErrorHandler = func(ctx *Context, err error) error {
	ctx.Response.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		ctx.Set("error", err)
		slog.Error("Handling HTTP route: " + err.Error())
	}
	return nil
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level plugs.
func (a *Kit) Add(method string, pattern string, handler Handler, plugs ...Plug) {
	a.routes = append(a.routes, Route{
		Path:   pattern,
		Method: method,
	})
	a.mux.HandleFunc(strings.Join([]string{method, pattern}, " "), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx := NewContext(w, r)
		h := handler
		for i := len(plugs) - 1; i >= 0; i-- {
			h = plugs[i](h)
		}
		for i := len(a.plugs) - 1; i >= 0; i-- {
			h = a.plugs[i](h)
		}
		if err := h(ctx); err != nil {
			if handleErr := a.ErrorHandler(ctx, err); handleErr != nil {
				slog.Error("Handling error: %v", handleErr)
			}
			return
		}
	})
}

// Connect registers a new CONNECT route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Connect(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodConnect, pattern, handler, plugs...)
}

// Delete registers a new DELETE route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Delete(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodDelete, pattern, handler, plugs...)
}

// Get registers a new GET route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Get(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodGet, pattern, handler, plugs...)
}

// Head registers a new HEAD route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Head(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodHead, pattern, handler, plugs...)
}

// Options registers a new OPTIONS route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Options(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodOptions, pattern, handler, plugs...)
}

// Post registers a new POST route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Post(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodPost, pattern, handler, plugs...)
}

// Put registers a new PUT route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Put(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodPut, pattern, handler, plugs...)
}

// Patch registers a new PATCH route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Patch(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodPatch, pattern, handler, plugs...)
}

// Trace registers a new TRACE route for a path with matching handler in the
// router with optional route-level plugs.
func (a *Kit) Trace(pattern string, handler Handler, plugs ...Plug) {
	a.Add(http.MethodTrace, pattern, handler, plugs...)
}
