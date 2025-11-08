package payload

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadPackageJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		wantErr bool
	}{
		"Valid package.json": {
			input: `{
				"name": "test-app",
				"dependencies": {
					"payload": "3.0.0",
					"@payloadcms/db-postgres": "3.0.0"
				},
				"devDependencies": {
					"typescript": "^5.0.0"
				}
			}`,
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{invalid json}`,
			wantErr: true,
		},
		"Empty file": {
			input:   `{}`,
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			path := "package.json"
			err := afero.WriteFile(fs, path, []byte(test.input), 0o644)
			require.NoError(t, err)

			pkg, err := ReadPackageJSON(fs, path)
			assert.Equal(t, test.wantErr, err != nil)

			if !test.wantErr {
				assert.NotNil(t, pkg)
				assert.NotNil(t, pkg.raw)
			}
		})
	}
}

func TestWritePackageJSON(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	path := "package.json"

	pkg := &PackageJSON{
		Dependencies: map[string]string{
			"payload": "3.0.0",
		},
		raw: map[string]interface{}{
			"name":    "test-app",
			"version": "1.0.0",
		},
	}

	err := WritePackageJSON(fs, path, pkg)
	require.NoError(t, err)

	// Verify file was written.
	exists, err := afero.Exists(fs, path)
	require.NoError(t, err)
	assert.True(t, exists)

	// Verify content is valid JSON.
	data, err := afero.ReadFile(fs, path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "test-app")
	assert.Contains(t, string(data), "payload")
}

func TestBumpPayloadDependencies(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg     *PackageJSON
		version string
		want    int
	}{
		"Bump main payload dependency": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
					"react":   "^18.0.0",
				},
			},
			version: "3.1.0",
			want:    1,
		},
		"Bump all payloadcms dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload":                      "3.0.0",
					"@payloadcms/db-postgres":      "3.0.0",
					"@payloadcms/richtext-lexical": "3.0.0",
					"react":                        "^18.0.0",
				},
			},
			version: "3.1.0",
			want:    3,
		},
		"Bump devDependencies": {
			pkg: &PackageJSON{
				DevDependencies: map[string]string{
					"payload": "3.0.0",
				},
			},
			version: "3.1.0",
			want:    1,
		},
		"Bump peerDependencies": {
			pkg: &PackageJSON{
				PeerDependencies: map[string]string{
					"payload": ">=3.0.0",
				},
			},
			version: "3.1.0",
			want:    1,
		},
		"No payload dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"react": "^18.0.0",
				},
			},
			version: "3.1.0",
			want:    0,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := BumpPayloadDependencies(test.pkg, test.version)
			assert.Len(t, result.Bumped, test.want)
			assert.Equal(t, test.version, result.NewVersion)

			// Verify dependencies were updated.
			if test.want > 0 {
				for _, dep := range result.Bumped {
					assert.NotEmpty(t, result.OldVersions[dep])
				}
			}
		})
	}
}

func TestIsPayloadDependency(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		name string
		want bool
	}{
		"Main payload package":       {name: "payload", want: true},
		"Scoped payloadcms package":  {name: "@payloadcms/db-postgres", want: true},
		"Another scoped package":     {name: "@payloadcms/richtext-lexical", want: true},
		"Non-payload package":        {name: "react", want: false},
		"Similar name but different": {name: "payloadjs", want: false},
		"Scoped but not payloadcms":  {name: "@types/node", want: false},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := isPayloadDependency(test.name)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestFormatVersion(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		version    string
		exactMatch bool
		want       string
	}{
		"Exact version":         {version: "3.0.0", exactMatch: true, want: "3.0.0"},
		"Caret version":         {version: "3.0.0", exactMatch: false, want: "^3.0.0"},
		"Version with v prefix": {version: "3.0.0", exactMatch: false, want: "^3.0.0"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := formatVersion(test.version, test.exactMatch)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestHasPayloadDependencies(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg  *PackageJSON
		want bool
	}{
		"Has payload in dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
				},
			},
			want: true,
		},
		"Has payloadcms scoped package": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"@payloadcms/db-postgres": "3.0.0",
				},
			},
			want: true,
		},
		"Has payload in devDependencies": {
			pkg: &PackageJSON{
				DevDependencies: map[string]string{
					"payload": "3.0.0",
				},
			},
			want: true,
		},
		"Has payload in peerDependencies": {
			pkg: &PackageJSON{
				PeerDependencies: map[string]string{
					"payload": ">=3.0.0",
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
		"Empty package": {
			pkg:  &PackageJSON{},
			want: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := HasPayloadDependencies(test.pkg)
			assert.Equal(t, test.want, got)
		})
	}
}
