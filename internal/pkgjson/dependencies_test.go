package pkgjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var payloadMatcher = func(name string) bool {
	return name == "payload" || strings.HasPrefix(name, "@payloadcms/")
}

func TestUpdateDependencies(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg  *PackageJSON
		want int
	}{
		"Updates matching dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
					"react":   "^18.0.0",
				},
			},
			want: 1,
		},
		"Updates across all dependency types": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
				},
				DevDependencies: map[string]string{
					"@payloadcms/eslint-config": "3.0.0",
				},
				PeerDependencies: map[string]string{
					"@payloadcms/db-postgres": "3.0.0",
				},
			},
			want: 3,
		},
		"No updates when no matches": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"react": "^18.0.0",
				},
			},
			want: 0,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := UpdateDependencies(test.pkg, payloadMatcher, func(name, version string) string {
				return "3.1.0"
			})

			assert.Len(t, result.Updated, test.want)
			if test.want > 0 {
				assert.NotEmpty(t, result.OldVersions)
			}
		})
	}
}

func TestPayloadMatcher(t *testing.T) {
	t.Parallel()

	matcher := func(name string) bool {
		return name == "payload" || strings.HasPrefix(name, "@payloadcms/")
	}

	assert.True(t, matcher("payload"))
	assert.True(t, matcher("@payloadcms/db-postgres"))
	assert.True(t, matcher("@payloadcms/richtext-lexical"))
	assert.False(t, matcher("react"))
	assert.False(t, matcher("payloadjs"))
	assert.False(t, matcher("@types/node"))
}

func TestFormatVersion(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		version  string
		useExact bool
		want     string
	}{
		"Exact version": {version: "3.0.0", useExact: true, want: "3.0.0"},
		"Caret version": {version: "3.0.0", useExact: false, want: "^3.0.0"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := FormatVersion(test.version, test.useExact)
			assert.Equal(t, test.want, got)
		})
	}
}
