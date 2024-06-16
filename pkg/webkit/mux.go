package webkit

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type (
	// Kit is the top-level framework instance for handling HTTP routes.
	//
	// Goroutine safety: Do not mutate WebKit instance fields after server has started. Accessing these
	// fields from handlers/middlewares and changing field values at the same time leads to data-races.
	// Adding new routes after the server has been started is also not safe!
	Kit struct {
		ErrorHandler    ErrorHandler
		NotFoundHandler Handler
		mux             *chi.Mux
		plugs           []Plug
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
		ErrorHandler:    DefaultErrorHandler,
		NotFoundHandler: DefaultNotFoundHandler,
		mux:             chi.NewRouter(),
		plugs:           []Plug{},
	}
}

// ServeHTTP implements the http.Handler interface.
func (k *Kit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k.mux.ServeHTTP(w, r)
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level plugs.
func (k *Kit) Add(method string, pattern string, handler Handler, plugs ...Plug) {
	k.mux.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		k.handle(w, r, handler, plugs...)
	})
}

// Plug adds a middleware function to the chain. These are called after
// the funcs that are passed directly to the route-level handlers.
//
// For example: app.Plug(middleware.Logger)
func (k *Kit) Plug(plugs ...Plug) {
	k.plugs = append(k.plugs, plugs...)
}

// Start starts the HTTP server.
func (k *Kit) Start(address string) error {
	server := &http.Server{
		Addr:    address,
		Handler: k.mux,
	}

	// Create a channel to receive signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		slog.Info("App listening on address: " + address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error: " + err.Error())
		}
	}()

	// Block until a signal is received
	<-sigChan

	// Log that shutdown signal is received
	slog.Info("Received shutdown signal. Shutting down...")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("HTTP shutdown error: " + err.Error())
	}

	return nil
}

// ErrorKey is the key used to store the error in the context.
const ErrorKey = "error"

// DefaultNotFoundHandler is the default handler that is called when no route
// matches the request.
var DefaultNotFoundHandler = func(ctx *Context) error {
	return ctx.String(http.StatusNotFound, "Not Found")
}

// DefaultErrorHandler is the default error handler that is called when a route
// handler returns an error.
var DefaultErrorHandler = func(ctx *Context, err error) error {
	return ctx.String(http.StatusInternalServerError, err.Error())
}

// PingHandler is a handler that returns a PONG response for health requests.
var PingHandler = func(ctx *Context) error {
	return ctx.String(http.StatusOK, "PONG")
}

// handle is a helper function that wraps the handler with plugs and executes
// the handler alongside any middleware functions.
func (k *Kit) handle(w http.ResponseWriter, r *http.Request, handler Handler, plugs ...Plug) {
	ctx := NewContext(w, r)

	h := handler
	for i := len(plugs) - 1; i >= 0; i-- {
		h = plugs[i](h)
	}

	for i := len(k.plugs) - 1; i >= 0; i-- {
		h = k.plugs[i](h)
	}

	if err := h(ctx); err != nil {
		k.handleError(ctx, err)
	}
}

func (k *Kit) handleError(ctx *Context, err error) {
	ctx.Set(ErrorKey, err)
	if handleErr := k.ErrorHandler(ctx, err); handleErr != nil {
		slog.Error("Handling error: %v", handleErr)
	}
}

// Connect registers a new CONNECT route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Connect(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodConnect, pattern, handler, plugs...)
}

// Delete registers a new DELETE route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Delete(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodDelete, pattern, handler, plugs...)
}

// Get registers a new GET route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Get(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodGet, pattern, handler, plugs...)
}

// Head registers a new HEAD route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Head(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodHead, pattern, handler, plugs...)
}

// Options registers a new OPTIONS route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Options(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodOptions, pattern, handler, plugs...)
}

// Post registers a new POST route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Post(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodPost, pattern, handler, plugs...)
}

// Put registers a new PUT route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Put(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodPut, pattern, handler, plugs...)
}

// Patch registers a new PATCH route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Patch(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodPatch, pattern, handler, plugs...)
}

// Trace registers a new TRACE route for a path with matching handler in the
// router with optional route-level plugs.
func (k *Kit) Trace(pattern string, handler Handler, plugs ...Plug) {
	k.Add(http.MethodTrace, pattern, handler, plugs...)
}

// NotFound sets a custom http.HandlerFunc for routing paths that could
// not be found. The default 404 handler is `http.NotFound`.
func (k *Kit) NotFound(handler Handler) {
	k.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		k.handle(w, r, handler)
	})
}
