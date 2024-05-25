package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseRecorder_Write(t *testing.T) {
	rr := NewResponseRecorder(httptest.NewRecorder())
	testBody := []byte("Test Body")

	_, err := rr.Write(testBody)
	require.NoError(t, err)

	assert.Equal(t, testBody, rr.Body.Bytes())
}

func TestResponseRecorder_WriteHeader(t *testing.T) {
	rr := NewResponseRecorder(httptest.NewRecorder())
	rr.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, rr.Status)
}
