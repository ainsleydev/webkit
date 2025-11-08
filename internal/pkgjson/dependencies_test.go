package pkgjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateDependencies(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg     *PackageJSON
		matcher DependencyMatcher
		want    int
	}{
		"Updates matching dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
					"react":   "^18.0.0",
				},
			},
			matcher: PayloadMatcher(),
			want:    1,
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
			matcher: PayloadMatcher(),
			want:    3,
		},
		"No updates when no matches": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"react": "^18.0.0",
				},
			},
			matcher: PayloadMatcher(),
			want:    0,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := UpdateDependencies(test.pkg, test.matcher, func(name, version string) string {
				return "3.1.0"
			})

			assert.Len(t, result.Updated, test.want)
			if test.want > 0 {
				assert.NotEmpty(t, result.OldVersions)
			}
		})
	}
}

func TestHasDependency(t *testing.T) {
	t.Parallel()

	pkg := &PackageJSON{
		Dependencies: map[string]string{
			"payload": "3.0.0",
		},
		DevDependencies: map[string]string{
			"typescript": "^5.0.0",
		},
	}

	assert.True(t, HasDependency(pkg, "payload"))
	assert.True(t, HasDependency(pkg, "typescript"))
	assert.False(t, HasDependency(pkg, "react"))
}

func TestHasAnyDependency(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg     *PackageJSON
		matcher DependencyMatcher
		want    bool
	}{
		"Has payload dependency": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
				},
			},
			matcher: PayloadMatcher(),
			want:    true,
		},
		"Has payloadcms scoped dependency": {
			pkg: &PackageJSON{
				DevDependencies: map[string]string{
					"@payloadcms/db-postgres": "3.0.0",
				},
			},
			matcher: PayloadMatcher(),
			want:    true,
		},
		"No payload dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"react": "^18.0.0",
				},
			},
			matcher: PayloadMatcher(),
			want:    false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := HasAnyDependency(test.pkg, test.matcher)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestIsDevDependency(t *testing.T) {
	t.Parallel()

	pkg := &PackageJSON{
		Dependencies: map[string]string{
			"payload": "3.0.0",
		},
		DevDependencies: map[string]string{
			"typescript": "^5.0.0",
		},
	}

	assert.True(t, IsDevDependency(pkg, "typescript"))
	assert.False(t, IsDevDependency(pkg, "payload"))
	assert.False(t, IsDevDependency(pkg, "react"))
}

func TestPayloadMatcher(t *testing.T) {
	t.Parallel()

	matcher := PayloadMatcher()

	assert.True(t, matcher("payload"))
	assert.True(t, matcher("@payloadcms/db-postgres"))
	assert.True(t, matcher("@payloadcms/richtext-lexical"))
	assert.False(t, matcher("react"))
	assert.False(t, matcher("payloadjs"))
	assert.False(t, matcher("@types/node"))
}

func TestSetMatcher(t *testing.T) {
	t.Parallel()

	deps := map[string]string{
		"payload": "3.0.0",
		"lexical": "0.28.0",
	}

	matcher := SetMatcher(deps)

	assert.True(t, matcher("payload"))
	assert.True(t, matcher("lexical"))
	assert.False(t, matcher("react"))
}

func TestFormatVersion(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		version   string
		useExact  bool
		want      string
	}{
		"Exact version":  {version: "3.0.0", useExact: true, want: "3.0.0"},
		"Caret version":  {version: "3.0.0", useExact: false, want: "^3.0.0"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := FormatVersion(test.version, test.useExact)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMergeDependencies(t *testing.T) {
	t.Parallel()

	t.Run("Adds new dependencies", func(t *testing.T) {
		t.Parallel()

		target := map[string]string{
			"react": "^18.0.0",
		}
		source := map[string]string{
			"payload": "3.0.0",
		}

		MergeDependencies(target, source, false)

		assert.Equal(t, "^18.0.0", target["react"])
		assert.Equal(t, "3.0.0", target["payload"])
	})

	t.Run("Preserves existing without overwrite", func(t *testing.T) {
		t.Parallel()

		target := map[string]string{
			"payload": "3.0.0",
		}
		source := map[string]string{
			"payload": "3.1.0",
		}

		MergeDependencies(target, source, false)

		assert.Equal(t, "3.0.0", target["payload"])
	})

	t.Run("Overwrites existing with overwrite", func(t *testing.T) {
		t.Parallel()

		target := map[string]string{
			"payload": "3.0.0",
		}
		source := map[string]string{
			"payload": "3.1.0",
		}

		MergeDependencies(target, source, true)

		assert.Equal(t, "3.1.0", target["payload"])
	})
}
