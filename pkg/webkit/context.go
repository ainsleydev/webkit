package webkit

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/goccy/go-json"
)

// Context represents the context of the current HTTP request. It holds request and
// response objects, path, path parameters, data and registered handler.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
}

// NewContext creates a new Context instance.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Response: w,
		Request:  r,
	}
}

// Get retrieves a value from the context using a key.
func (c *Context) Get(key string) any {
	return c.Request.Context().Value(key)
}

// Set sets a value into the context using a key.
func (c *Context) Set(key string, value any) {
	ctx := context.WithValue(c.Request.Context(), key, value)
	c.Request = c.Request.WithContext(ctx)
}

// Context returns the original request context.
func (c *Context) Context() context.Context {
	return c.Request.Context()
}

// Param retrieves a parameter from the route parameters.
func (c *Context) Param(key string) string {
	return c.Request.PathValue(key) // TODO, need to use chi
}

// Render renders a templated component to the response writer.
func (c *Context) Render(component templ.Component) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(http.StatusOK)
	return component.Render(c.Context(), c.Response)
}

// RenderWithStatus renders a templated component to the response writer with the
// specified status code.
func (c *Context) RenderWithStatus(status int, component templ.Component) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(status)
	return component.Render(c.Context(), c.Response)
}

// Redirect performs an HTTP redirect to the specified URL with the given status code.
func (c *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return errors.New("invalid redirect code")
	}
	c.Response.Header().Set("Location", url)
	c.Response.WriteHeader(code)
	return nil
}

// NoContent sends a response with no response body and the provided status code.
func (c *Context) NoContent(status int) error {
	c.Response.WriteHeader(status)
	return nil
}

// String writes a plain text response with the provided status code and data.
// The header is set to text/plain.
func (c *Context) String(status int, v string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(v))
	return err
}

// JSON writes a JSON response with the provided status code and data.
// The header is set to application/json.
func (c *Context) JSON(status int, v any) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	return json.NewEncoder(c.Response).Encode(v)
}

// HTML writes an HTML response with the provided status code and data.
// The header is set to text/html; charset=utf-8
// TODO: Unit test this func.
func (c *Context) HTML(status int, v string) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(v))
	return err
}

// BindJSON decodes the JSON request body and stores the result
// in the value pointed to by v.
func (c *Context) BindJSON(v any) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (c *Context) IsTLS() bool {
	return c.Request.TLS != nil
}

// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
func (c *Context) IsWebSocket() bool {
	upgrade := c.Request.Header.Get("Upgrade")
	return strings.EqualFold(upgrade, "websocket")
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (c *Context) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.IsTLS() {
		return "https"
	}
	if scheme := c.Request.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if scheme := c.Request.Header.Get("X-Forwarded-Protocol"); scheme != "" {
		return scheme
	}
	if ssl := c.Request.Header.Get("X-Forwarded-Ssl"); ssl == "on" {
		return "https"
	}
	if scheme := c.Request.Header.Get("X-Url-Scheme"); scheme != "" {
		return scheme
	}
	return "http"
}
