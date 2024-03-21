package webkit

import (
	"context"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/goccy/go-json"
)

// Context represents the context of the current HTTP request. It holds request and
// response objects, path, path parameters, data and registered handler.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	ctx      context.Context
	webKit   *Kit
}

// NewContext creates a new Context instance.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Response: w,
		Request:  r,
		ctx:      r.Context(),
	}
}

// Get retrieves a value from the context using a key.
func (c *Context) Get(key string) any {
	return c.ctx.Value(key)
}

// Set sets a value into the context using a key.
func (c *Context) Set(key string, value any) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

// Param retrieves a parameter from the route parameters.
func (c *Context) Param(key string) string {
	return c.Request.PathValue(key)
}

// Render renders a templated component to the response writer.
func (c *Context) Render(component templ.Component) error {
	return component.Render(c.ctx, c.Response)
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

// JSON writes a JSON response with the provided status code and data.
// The header is set to application/json.
func (c *Context) JSON(status int, v any) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	return json.NewEncoder(c.Response).Encode(v)
}

// String writes a plain text response with the provided status code and data.
// The header is set to text/plain.
func (c *Context) String(status int, v string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(v))
	return err
}
