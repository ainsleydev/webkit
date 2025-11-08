package pkgjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasDependency(t *testing.T) {
	t.Parallel()

	pkg := &PackageJSON{
		Dependencies: map[string]string{
			"react": "18.0.0",
		},
		DevDependencies: map[string]string{
			"typescript": "^5.0.0",
		},
		PeerDependencies: map[string]string{
			"graphql": "^16.0.0",
		},
	}

	assert.True(t, pkg.HasDependency("react"))
	assert.True(t, pkg.HasDependency("typescript"))
	assert.True(t, pkg.HasDependency("graphql"))
	assert.False(t, pkg.HasDependency("vue"))
}

func TestHasAnyDependency(t *testing.T) {
	t.Parallel()

	containsMatcher := func(substr string) DependencyMatcher {
		return func(name string) bool {
			return strings.Contains(name, substr)
		}
	}

	tt := map[string]struct {
		pkg     *PackageJSON
		matcher DependencyMatcher
		want    bool
	}{
		"Matches in Dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"react": "^18.0.0",
					"vue":   "^3.0.0",
				},
			},
			matcher: containsMatcher("react"),
			want:    true,
		},
		"Matches in DevDependencies": {
			pkg: &PackageJSON{
				DevDependencies: map[string]string{
					"eslint": "^8.0.0",
				},
			},
			matcher: containsMatcher("eslint"),
			want:    true,
		},
		"Matches in PeerDependencies": {
			pkg: &PackageJSON{
				PeerDependencies: map[string]string{
					"typescript": "^5.0.0",
				},
			},
			matcher: containsMatcher("typescript"),
			want:    true,
		},
		"No matching dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"axios": "^1.0.0",
				},
			},
			matcher: containsMatcher("react"),
			want:    false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.pkg.HasAnyDependency(test.matcher)
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

	assert.True(t, pkg.IsDevDependency("typescript"))
	assert.False(t, pkg.IsDevDependency("payload"))
	assert.False(t, pkg.IsDevDependency("react"))
}
