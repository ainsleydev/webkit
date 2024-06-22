package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

var noContextHandler = func(ctx *webkit.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func TestAddTrailingSlash(t *testing.T) {
	t.Parallel()

	host := "test.com"

	tt := map[string]struct {
		url          string
		wantLocation string
		wantStatus   int
	}{
		"Home": {
			url:          "/",
			wantLocation: "",
			wantStatus:   http.StatusOK,
		},
		"With Redirect": {
			url:          "/test",
			wantLocation: "//" + host + "/test/",
			wantStatus:   http.StatusMovedPermanently,
		},
		"With Query": {
			url:          "/test?query=test",
			wantLocation: "//" + host + "/test/?query=test",
			wantStatus:   http.StatusMovedPermanently,
		},
		"No Redirect": {
			url:          "/test/",
			wantLocation: "",
			wantStatus:   http.StatusOK,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			app := webkit.New()
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, test.url, nil)
			req.Host = host
			require.NoError(t, err)

			app.Plug(TrailingSlashRedirect)
			app.Get("/", noContextHandler)
			app.Get("/test", noContextHandler)
			app.Get("/test/", noContextHandler)
			app.ServeHTTP(rr, req)

			assert.Equal(t, test.wantStatus, rr.Code)
			assert.Equal(t, test.wantLocation, rr.Header().Get("Location"))
		})
	}
}

func TestWWWRedirect(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		host         string
		wantLocation string
		wantStatus   int
	}{
		"With Redirect": {
			host:         "example.com",
			wantLocation: "https://www.example.com/",
			wantStatus:   http.StatusMovedPermanently,
		},
		"No Redirect": {
			host:         "www.example.com",
			wantLocation: "",
			wantStatus:   http.StatusOK,
		},
		"Subdomain": {
			host:         "sub.example.com",
			wantLocation: "https://www.sub.example.com/",
			wantStatus:   http.StatusMovedPermanently,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			rr := redirectTest(t, WWWRedirect, tc.host)
			assert.Equal(t, tc.wantStatus, rr.Code)
			assert.Equal(t, tc.wantLocation, rr.Header().Get("Location"))
		})
	}
}

func TestNonWWWRedirect(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		host         string
		wantLocation string
		wantStatus   int
	}{
		"With Redirect": {
			host:         "www.example.com",
			wantLocation: "https://example.com/",
			wantStatus:   http.StatusMovedPermanently,
		},
		"No Redirect": {
			host:         "example.com",
			wantLocation: "",
			wantStatus:   http.StatusOK,
		},
		"Subdomain": {
			host:         "sub.example.com",
			wantLocation: "",
			wantStatus:   http.StatusOK,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			rr := redirectTest(t, NonWWWRedirect, tc.host)
			assert.Equal(t, tc.wantStatus, rr.Code)
			assert.Equal(t, tc.wantLocation, rr.Header().Get("Location"))
		})
	}
}

func redirectTest(t *testing.T, middleware webkit.Plug, host string) *httptest.ResponseRecorder {
	t.Helper()

	app := webkit.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = host
	rr := httptest.NewRecorder()

	app.Get("/", noContextHandler, middleware)
	app.ServeHTTP(rr, req)

	return rr
}
