package payloadcms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultBody    = []byte(`{"id": 1, "name": "John Doe"}`)
	defaultHandler = func(t *testing.T) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(defaultBody))
			assert.NoError(t, err)
		})
	}
)

func Setup(t *testing.T, handlerFunc http.HandlerFunc, baseURL string) (*Client, func()) {
	t.Helper()

	server := httptest.NewServer(handlerFunc)
	return &Client{
			baseURL: server.URL,
			client:  server.Client(),
		}, func() {
			server.Close()
		}
}

func TestClientDo(t *testing.T) {

	tcs := map[string]struct {
		method   string
		path     string
		wantCode int
		wantBody []byte
		wantErr  bool
	}{
		"Happy Case - 200 OK": {
			method:   http.MethodGet,
			path:     "/users/1",
			wantCode: http.StatusOK,
			wantBody: defaultBody,
			wantErr:  false,
		},
		"Error Case - Request Creation": {
			method:  "INVALID",
			path:    "werong",
			wantErr: true,
		},
	}

	for name, test := range tcs {
		t.Run(name, func(t *testing.T) {
			client, teardown := Setup(t, defaultHandler(t), string(test.wantBody))
			defer teardown()

			response, err := client.Do(context.Background(), test.method, test.path, nil, nil)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantCode, response.StatusCode)
			assert.Equal(t, test.wantBody, response.Content)
		})
	}
}
