package payload

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestShouldSkipMiddleware(t *testing.T) {
	tt := map[string]struct {
		input *webkit.Context
		want  bool
	}{
		"File request should be skipped": {
			input: &webkit.Context{
				Request: &http.Request{
					URL:    &url.URL{Path: "/static/file.txt"},
					Method: http.MethodGet,
				},
			},
			want: true,
		},
		"Non-GET request should be skipped": {
			input: &webkit.Context{
				Request: &http.Request{
					URL:    &url.URL{Path: "/"},
					Method: http.MethodPost,
				},
			},
			want: true,
		},
		"Skippable path should be skipped": {
			input: &webkit.Context{
				Request: &http.Request{
					URL:    &url.URL{Path: skippable[0]},
					Method: http.MethodGet,
				},
			},
			want: true,
		},
		"Passed": {
			input: &webkit.Context{
				Request: &http.Request{
					URL:    &url.URL{Path: "/home"},
					Method: http.MethodGet,
				},
			},
			want: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := shouldSkipMiddleware(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
