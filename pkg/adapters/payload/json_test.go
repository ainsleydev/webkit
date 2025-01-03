package payload

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestJSON_MarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   JSON
		want    string
		wantErr bool
	}{
		"Simple JSON Object": {
			input:   JSON{"key": "value"},
			want:    `{"key":"value"}`,
			wantErr: false,
		},
		"Empty JSON Object": {
			input:   JSON{},
			want:    "{}",
			wantErr: false,
		},
		"Marshal Error (Unsupported Type)": {
			input:   JSON{"key": func() {}}, // Functions are unsupported in JSON
			want:    "",
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := test.input.MarshalJSON()

			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, string(data))
		})
	}
}

func TestJSON_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		want    JSON
		wantErr bool
	}{
		"Valid JSON Object": {
			input:   `{"key":"value","another_key":42}`,
			want:    JSON{"key": "value", "another_key": float64(42)}, // Unmarshal converts numbers to float64 by default
			wantErr: false,
		},
		"Empty JSON Object": {
			input:   `{}`,
			want:    JSON{},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{"key": "value"`,
			want:    nil,
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var j JSON
			err := json.Unmarshal([]byte(test.input), &j)

			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, j)
		})
	}
}
