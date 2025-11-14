package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolSpec_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		json    string
		want    ToolSpec
		wantErr bool
	}{
		"String version": {
			json: `"v0.2.543"`,
			want: ToolSpec{Version: "v0.2.543"},
		},
		"String latest": {
			json: `"latest"`,
			want: ToolSpec{Version: "latest"},
		},
		"String disabled": {
			json: `"disabled"`,
			want: ToolSpec{Version: "disabled"},
		},
		"Empty string": {
			json: `""`,
			want: ToolSpec{Version: "disabled"},
		},
		"Object with version only": {
			json: `{"version": "v1.0.0"}`,
			want: ToolSpec{Version: "v1.0.0"},
		},
		"Object with version and install": {
			json: `{"version": "v1.0.0", "install": "custom install command"}`,
			want: ToolSpec{Version: "v1.0.0", Install: "custom install command"},
		},
		"Object with install only": {
			json: `{"install": "curl -sSL https://example.com | sh"}`,
			want: ToolSpec{Install: "curl -sSL https://example.com | sh"},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got ToolSpec
			err := got.UnmarshalJSON([]byte(test.json))

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestReplaceVersion(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		install    string
		oldVersion string
		newVersion string
		want       string
	}{
		"Go install latest to specific": {
			install:    "go install github.com/a-h/templ/cmd/templ@latest",
			oldVersion: "latest",
			newVersion: "v0.2.543",
			want:       "go install github.com/a-h/templ/cmd/templ@v0.2.543",
		},
		"Go install specific to specific": {
			install:    "go install github.com/a-h/templ/cmd/templ@v0.2.500",
			oldVersion: "v0.2.500",
			newVersion: "v0.2.543",
			want:       "go install github.com/a-h/templ/cmd/templ@v0.2.543",
		},
		"Version appears multiple times": {
			install:    "go install github.com/example/tool@latest && tool@latest --version",
			oldVersion: "latest",
			newVersion: "v1.0.0",
			want:       "go install github.com/example/tool@v1.0.0 && tool@v1.0.0 --version",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := replaceVersion(test.install, test.oldVersion, test.newVersion)
			assert.Equal(t, test.want, got)
		})
	}
}
