package wordpress

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Setup(t *testing.T, serverResponse string, serverStatus int) (*Client, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(serverStatus)
		_, err := io.WriteString(w, serverResponse)
		require.NoError(t, err)
	}))

	opts := NewOptions().WithBaseURL(server.URL)
	client, err := New(opts)
	require.NoError(t, err)

	teardown := func() {
		server.Close()
	}

	return client, teardown
}

func TestClient_Get(t *testing.T) {
	tt := map[string]struct {
		url            string
		serverResponse string
		serverStatus   int
		wantBody       []byte
		wantErr        bool
	}{
		"Valid Request": {
			url:            "/test",
			serverResponse: "test response",
			serverStatus:   http.StatusOK,
			wantBody:       []byte("test response"),
			wantErr:        false,
		},
		"Bad URL": {
			url:          "£$%£$/invalid",
			serverStatus: http.StatusNotFound,
			wantBody:     nil,
			wantErr:      true,
		},
		"Invalid Request": {
			url:          "/invalid",
			serverStatus: http.StatusNotFound,
			wantBody:     nil,
			wantErr:      true,
		},
		"Non 200 Status": {
			url:          "/non200",
			serverStatus: http.StatusInternalServerError,
			wantBody:     nil,
			wantErr:      true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			client, teardown := Setup(t, test.serverResponse, test.serverStatus)
			defer teardown()

			got, err := client.Get(context.TODO(), test.url)
			assert.Equal(t, test.wantBody, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestClient_Get_ReadError(t *testing.T) {
	client, teardown := Setup(t, "test response", http.StatusOK)
	defer teardown()

	client.reader = func(r io.Reader) ([]byte, error) {
		return nil, assert.AnError
	}

	_, err := client.Get(context.TODO(), "/test")
	assert.Error(t, err)
}

func TestClient_GetAndUnmarshal(t *testing.T) {
	type TestStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tt := map[string]struct {
		url            string
		serverResponse string
		serverStatus   int
		wantStruct     TestStruct
		wantErr        bool
	}{
		"Valid Request": {
			url:            "/test",
			serverResponse: `{"id": 1, "name": "test"}`,
			serverStatus:   http.StatusOK,
			wantStruct:     TestStruct{ID: 1, Name: "test"},
			wantErr:        false,
		},
		"Invalid JSON": {
			url:            "/invalid",
			serverStatus:   http.StatusOK,
			serverResponse: `{WRONG{{`,
			wantErr:        true,
		},
		"Non 200 Status": {
			url:          "/non200",
			serverStatus: http.StatusInternalServerError,
			wantStruct:   TestStruct{},
			wantErr:      true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			client, teardown := Setup(t, test.serverResponse, test.serverStatus)
			defer teardown()

			var result TestStruct
			err := client.GetAndUnmarshal(context.TODO(), test.url, &result)

			assert.Equal(t, test.wantStruct, result)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
