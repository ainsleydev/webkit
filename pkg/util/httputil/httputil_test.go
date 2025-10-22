package httputil

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs2xx(t *testing.T) {
	tt := map[string]struct {
		input int
		want  bool
	}{
		"200": {
			input: 200,
			want:  true,
		},
		"204": {
			input: 204,
			want:  true,
		},
		"300": {
			input: 300,
			want:  false,
		},
		"199": {
			input: 199,
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := Is2xx(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestIs3xx(t *testing.T) {
	tt := map[string]struct {
		input int
		want  bool
	}{
		"300": {
			input: 300,
			want:  true,
		},
		"301": {
			input: 301,
			want:  true,
		},
		"400": {
			input: 400,
			want:  false,
		},
		"299": {
			input: 299,
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := Is3xx(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestIs4xx(t *testing.T) {
	tt := map[string]struct {
		input int
		want  bool
	}{
		"400": {
			input: 400,
			want:  true,
		},
		"401": {
			input: 401,
			want:  true,
		},
		"500": {
			input: 500,
			want:  false,
		},
		"399": {
			input: 399,
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := Is4xx(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestIs5xx(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input int
		want  bool
	}{
		"500": {
			input: 500,
			want:  true,
		},
		"501": {
			input: 501,
			want:  true,
		},
		"600": {
			input: 600,
			want:  false,
		},
		"499": {
			input: 499,
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := Is5xx(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestIsFileRequest(t *testing.T) {
	tt := map[string]struct {
		input string
		want  bool
	}{
		"CSS file": {
			input: "/styles/main.css",
			want:  true,
		},
		"JS file": {
			input: "/scripts/app.js",
			want:  true,
		},
		"Root": {
			input: "/",
			want:  false,
		},
		"Page Path": {
			input: "/about",
			want:  false,
		},
		"Trailing slash": {
			input: "/about/team",
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					Path: test.input,
				},
			}
			got := IsFileRequest(req)
			assert.Equal(t, test.want, got)
		})
	}
}
