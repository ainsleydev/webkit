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

func TestResponseRecorder_Write_DefaultsStatusTo200(t *testing.T) {
	// Write without a prior WriteHeader must record 200, mirroring net/http behaviour.
	rr := NewResponseRecorder(httptest.NewRecorder())
	_, err := rr.Write([]byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Status)
}

func TestResponseRecorder_Write_PreservesExplicitStatus(t *testing.T) {
	rr := NewResponseRecorder(httptest.NewRecorder())
	rr.WriteHeader(http.StatusCreated)
	_, err := rr.Write([]byte("created"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rr.Status)
}

func TestResponseRecorder_WriteHeader(t *testing.T) {
	rr := NewResponseRecorder(httptest.NewRecorder())
	rr.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, rr.Status)
}
