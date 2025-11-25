package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config Config
		key    string
		want   string
		wantOk bool
	}{
		"Returns String When Key Exists": {
			config: Config{"name": "test"},
			key:    "name",
			want:   "test",
			wantOk: true,
		},
		"Returns Empty String When Key Missing": {
			config: Config{"name": "test"},
			key:    "missing",
			want:   "",
			wantOk: false,
		},
		"Returns Empty String When Value Not String": {
			config: Config{"count": 42},
			key:    "count",
			want:   "",
			wantOk: false,
		},
		"Returns Empty String When Config Nil": {
			config: nil,
			key:    "name",
			want:   "",
			wantOk: false,
		},
		"Returns Empty String When Config Empty": {
			config: Config{},
			key:    "name",
			want:   "",
			wantOk: false,
		},
		"Handles Empty String Value": {
			config: Config{"name": ""},
			key:    "name",
			want:   "",
			wantOk: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, ok := tt.config.String(tt.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestConfig_Int(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config Config
		key    string
		want   int
		wantOk bool
	}{
		"Returns Int When Key Exists": {
			config: Config{"count": 42},
			key:    "count",
			want:   42,
			wantOk: true,
		},
		"Returns Zero When Key Missing": {
			config: Config{"count": 42},
			key:    "missing",
			want:   0,
			wantOk: false,
		},
		"Returns Zero When Value Not Int": {
			config: Config{"name": "test"},
			key:    "name",
			want:   0,
			wantOk: false,
		},
		"Returns Zero When Config Nil": {
			config: nil,
			key:    "count",
			want:   0,
			wantOk: false,
		},
		"Returns Zero When Config Empty": {
			config: Config{},
			key:    "count",
			want:   0,
			wantOk: false,
		},
		"Handles Zero Value": {
			config: Config{"count": 0},
			key:    "count",
			want:   0,
			wantOk: true,
		},
		"Handles Negative Values": {
			config: Config{"count": -10},
			key:    "count",
			want:   -10,
			wantOk: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, ok := tt.config.Int(tt.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestConfig_Bool(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config Config
		key    string
		want   bool
		wantOk bool
	}{
		"Returns True When Key Exists": {
			config: Config{"enabled": true},
			key:    "enabled",
			want:   true,
			wantOk: true,
		},
		"Returns False When Key Exists": {
			config: Config{"enabled": false},
			key:    "enabled",
			want:   false,
			wantOk: true,
		},
		"Returns False When Key Missing": {
			config: Config{"enabled": true},
			key:    "missing",
			want:   false,
			wantOk: false,
		},
		"Returns False When Value Not Bool": {
			config: Config{"name": "test"},
			key:    "name",
			want:   false,
			wantOk: false,
		},
		"Returns False When Config Nil": {
			config: nil,
			key:    "enabled",
			want:   false,
			wantOk: false,
		},
		"Returns False When Config Empty": {
			config: Config{},
			key:    "enabled",
			want:   false,
			wantOk: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, ok := tt.config.Bool(tt.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}
