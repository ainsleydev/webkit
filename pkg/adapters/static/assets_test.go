package static

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveAssetPath(t *testing.T) {
	tt := map[string]struct {
		input string
		want  string
	}{
		"Basic asset path": {
			input: "/assets/image.jpg",
			want:  "dist/image.jpg",
		},
		"Multiple subdirectories": {
			input: "/assets/images/subfolder/file.png",
			want:  "dist/images/subfolder/file.png",
		},
		"Root path": {
			input: "/",
			want:  "/",
		},
		"Single file without subdirectory": {
			input: "file.txt",
			want:  "file.txt",
		},
		"Empty string": {
			input: "",
			want:  "",
		},
		"Path without leading slash": {
			input: "assets/script.js",
			want:  "dist/script.js",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := AssetToBasePath(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
