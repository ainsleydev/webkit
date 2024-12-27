package payload

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPoint_MarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Point
		want  string
	}{
		"Los Angeles": {
			input: Point{Longitude: -118.243683, Latitude: 34.052235},
			want:  "[-118.243683,34.052235]",
		},
		"London": {
			input: Point{Longitude: -0.1278, Latitude: 51.5074},
			want:  "[-0.1278,51.5074]",
		},
		"Zero Point": {
			input: Point{Longitude: 0, Latitude: 0},
			want:  "[0,0]",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := json.Marshal(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestPoint_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  Point
		err   bool
	}{
		"Valid Point Los Angeles": {
			input: "[-118.243683,34.052235]",
			want:  Point{Longitude: -118.243683, Latitude: 34.052235},
		},
		"Valid Point London": {
			input: "[-0.1278,51.5074]",
			want:  Point{Longitude: -0.1278, Latitude: 51.5074},
		},
		"Zero Point": {
			input: "[0,0]",
			want:  Point{Longitude: 0, Latitude: 0},
		},
		"Invalid JSON": {
			input: "not json",
			err:   true,
		},
		"Wrong Number of Elements": {
			input: "[-118.243683]",
			err:   true,
		},
		"Too Many Elements": {
			input: "[-118.243683,34.052235,0]",
			err:   true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var got Point
			err := json.Unmarshal([]byte(test.input), &got)
			if test.err {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestPoint_ToSlice(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Point
		want  []float64
	}{
		"Los Angeles": {
			input: Point{Longitude: -118.243683, Latitude: 34.052235},
			want:  []float64{-118.243683, 34.052235},
		},
		"London": {
			input: Point{Longitude: -0.1278, Latitude: 51.5074},
			want:  []float64{-0.1278, 51.5074},
		},
		"Zero Point": {
			input: Point{},
			want:  []float64{0, 0},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.ToSlice()
			assert.Equal(t, test.want, got)
		})
	}
}
