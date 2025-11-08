package pkgjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	assert.True(t, pkg.HasDependency("payload"))
	assert.True(t, pkg.HasDependency("typescript"))
	assert.False(t, pkg.HasDependency("react"))
}

func TestHasAnyDependency(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg  *PackageJSON
		want bool
	}{
		"Has payload dependency": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
				},
			},
			want: true,
		},
		"Has payloadcms scoped dependency": {
			pkg: &PackageJSON{
				DevDependencies: map[string]string{
					"@payloadcms/db-postgres": "3.0.0",
				},
			},
			want: true,
		},
		"No payload dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"react": "^18.0.0",
				},
			},
			want: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.pkg.HasAnyDependency(payloadMatcher)
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
