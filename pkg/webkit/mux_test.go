package webkit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var handler = func(ctx *Context) error {
	return ctx.String(http.StatusOK, "test")
}

func HandlerTest(t *testing.T, kit *Kit) {
	t.Helper()
	app := New()
	app.Get("/", func(ctx *Context) error {
		return ctx.String(http.StatusOK, "test")
	})
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, rr.Body.String(), "test")
}

func TestAdd(t *testing.T) {
	t.Run("Valid request", func(t *testing.T) {
		app := New()
		app.Add(http.MethodGet, "/", handler)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		app.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Not allowed", func(t *testing.T) {
		app := New()
		app.Add(http.MethodPost, "/", handler)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		app.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})
}

func TestKit_Connect(t *testing.T) {
	app := New()
	app.Connect("/", handler)
	HandlerTest(t, app)
}

func TestKit_Delete(t *testing.T) {
	app := New()
	app.Delete("/", handler)
	HandlerTest(t, app)
}

func TestKit_Get(t *testing.T) {
	app := New()
	app.Get("/", handler)
	HandlerTest(t, app)
}

func TestKit_Head(t *testing.T) {
	app := New()
	app.Head("/", handler)
	HandlerTest(t, app)
}

func TestKit_Options(t *testing.T) {
	app := New()
	app.Options("/", handler)
	HandlerTest(t, app)
}

func TestKit_Post(t *testing.T) {
	app := New()
	app.Post("/", handler)
	HandlerTest(t, app)
}

func TestKit_Put(t *testing.T) {
	app := New()
	app.Put("/", handler)
	HandlerTest(t, app)
}

func TestKit_Patch(t *testing.T) {
	app := New()
	app.Patch("/", handler)
	HandlerTest(t, app)
}

func TestKit_Trace(t *testing.T) {
	app := New()
	app.Trace("/", handler)
	HandlerTest(t, app)
}
