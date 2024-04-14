package wordpress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_Validate(t *testing.T) {
	tests := map[string]struct {
		options Options
		wantErr bool
	}{
		"Valid": {
			options: Options{baseURL: "https://example.com"},
			wantErr: false,
		},
		"Empty BaseURL": {
			options: Options{},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.options.Validate()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestOptions_WithBaseURL(t *testing.T) {
	baseURL := "https://example.com"
	options := Options{}.WithBaseURL(baseURL)
	assert.Equal(t, options.baseURL, baseURL)
}

func TestOptions_WithBasicAuth(t *testing.T) {
	user := "user"
	password := "password"
	options := Options{}.WithBasicAuth(user, password)
	assert.Equal(t, options.hasBasicAuth, true)
	assert.Equal(t, options.authUser, user)
	assert.Equal(t, options.authPassword, password)
}
