package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint_Latitude(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Point
		want  float64
	}{
		"Valid Point - Latitude":         {input: Point{-118.243683, 34.052235}, want: 34.052235},
		"Another Valid Point - Latitude": {input: Point{-0.1278, 51.5074}, want: 51.5074},
		"Empty Point":                    {input: Point{}, want: 0},
		"Only Longitude in Point":        {input: Point{-118.243683}, want: 0},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.Latitude()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestPoint_Longitude(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Point
		want  float64
	}{
		"Valid Point - Longitude":         {input: Point{-118.243683, 34.052235}, want: -118.243683},
		"Another Valid Point - Longitude": {input: Point{-0.1278, 51.5074}, want: -0.1278},
		"Empty Point":                     {input: Point{}, want: 0},
		"Only Latitude in Point":          {input: Point{34.052235}, want: 34.052235},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.Longitude()
			assert.Equal(t, test.want, got)
		})
	}
}
