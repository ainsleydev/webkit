package payload

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestNew_Error(t *testing.T) {
	tt := map[string]struct {
		options []Option
		envURL  string
		want    string
	}{
		"Failed Validation": {
			options: []Option{},
			envURL:  "https://api.payloadcms.com",
			want:    "kit is required",
		},
		"Client error": {
			options: []Option{
				WithWebkit(webkit.New()),
			},
			envURL: "https://api.payloadcms.com",
			want:   "baseURL is required",
		},
		"Set environment variable error": {
			options: []Option{
				WithBaseURL(string([]byte{0x7f})),
				WithAPIKey("test-api-key"),
				WithWebkit(webkit.New()),
			},
			envURL: "",
			want:   "invalid argument",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			orig := envPayloadURL
			defer func() { envPayloadURL = orig }()
			envPayloadURL = test.envURL

			_, err := New(test.options...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.want)
		})
	}
}

func TestNew_OK(t *testing.T) {
	var (
		baseURL = "https://api.payloadcms.com"
		apiKey  = "api-key"
	)

	t.Setenv(env.AppEnvironmentKey, env.Production)

	got, err := New(
		WithWebkit(webkit.New()),
		WithCache(cache.NewInMemory(time.Hour*24)),
		WithBaseURL(baseURL),
		WithAPIKey(apiKey),
		WithNavigation(),
		WithMaintenanceHandler(defaultMaintenanceRenderer),
		WithGlobalMiddleware[string]("navigation"),
		WithSilentLogs(),
	)

	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, got.baseURL, baseURL)
	assert.Equal(t, got.apiKey, apiKey)
	assert.Equal(t, got.env[env.AppEnvironmentKey], env.Production)
	assert.Equal(t, os.Getenv(envPayloadURL), baseURL)
}
