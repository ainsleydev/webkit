package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceProject(t *testing.T) {
	t.Parallel()

	got := SourceProject()
	assert.Equal(t, "project", got)
}

func TestSourceApp(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple App Name": {
			input: "web",
			want:  "app:web",
		},
		"App Name With Dash": {
			input: "api-server",
			want:  "app:api-server",
		},
		"Empty App Name": {
			input: "",
			want:  "app:",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := SourceApp(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResource(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple Resource Name": {
			input: "postgres",
			want:  "resource:postgres",
		},
		"Resource Name With Dash": {
			input: "redis-cache",
			want:  "resource:redis-cache",
		},
		"Empty Resource Name": {
			input: "",
			want:  "resource:",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := SourceResource(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
