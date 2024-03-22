package webkit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapHandlerFunc(t *testing.T) {
	called := false
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	wrappedHandler := WrapHandlerFunc(handlerFunc)
	err := wrappedHandler(&Context{
		Response: httptest.NewRecorder(),
		Request:  httptest.NewRequest(http.MethodGet, "/", nil),
	})
	require.NoError(t, err)
	assert.True(t, called)
}

func TestWrapHandler(t *testing.T) {
	called := false
	wrappedHandler := WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	err := wrappedHandler(&Context{
		Response: httptest.NewRecorder(),
		Request:  httptest.NewRequest(http.MethodGet, "/", nil),
	})
	require.NoError(t, err)
	assert.True(t, called)
}
