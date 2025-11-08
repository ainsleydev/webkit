package pkgjson

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		wantErr bool
	}{
		"Valid package.json": {
			input: `{
				"name": "test-app",
				"version": "1.0.0",
				"dependencies": {
					"payload": "3.0.0",
					"react": "^18.0.0"
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

			pkg, err := Read(fs, path)
			assert.Equal(t, test.wantErr, err != nil)

			if !test.wantErr {
				assert.NotNil(t, pkg)
				assert.NotNil(t, pkg.raw)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	path := "package.json"

	pkg := &PackageJSON{
		Name:    "test-app",
		Version: "1.0.0",
		Dependencies: map[string]string{
			"payload": "3.0.0",
		},
		raw: map[string]interface{}{
			"name":    "test-app",
			"version": "1.0.0",
		},
	}

	err := Write(fs, path, pkg)
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

func TestExists(t *testing.T) {
	t.Parallel()

	t.Run("File exists", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"

		err := afero.WriteFile(fs, path, []byte(`{}`), 0o644)
		require.NoError(t, err)

		exists, err := Exists(fs, path)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		exists, err := Exists(fs, "package.json")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestReadWrite(t *testing.T) {
	t.Parallel()

	t.Run("Preserves all fields", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"

		original := `{
	"name": "test-app",
	"version": "1.0.0",
	"customField": "custom-value",
	"dependencies": {
		"payload": "3.0.0"
	}
}`
		err := afero.WriteFile(fs, path, []byte(original), 0o644)
		require.NoError(t, err)

		// Read the package.json.
		pkg, err := Read(fs, path)
		require.NoError(t, err)

		// Modify a dependency.
		pkg.Dependencies["payload"] = "3.1.0"

		// Write it back.
		err = Write(fs, path, pkg)
		require.NoError(t, err)

		// Read it again.
		data, err := afero.ReadFile(fs, path)
		require.NoError(t, err)

		// Verify custom field is preserved.
		assert.Contains(t, string(data), "custom-value")
		assert.Contains(t, string(data), "3.1.0")
	})
}
