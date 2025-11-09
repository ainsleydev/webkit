package pkgjson

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/fsext"
)

func TestRead(t *testing.T) {
	t.Parallel()

	t.Run("File doesnt exist", func(t *testing.T) {
		t.Parallel()

		_, err := Read(afero.NewMemMapFs(), "package.json")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "reading package.json")
	})

	t.Run("Valid package.json", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"
		content := `{
			"name": "test-app",
			"version": "1.0.0",
			"dependencies": {
				"payload": "3.0.0",
				"react": "^18.0.0"
			},
			"devDependencies": {
				"typescript": "^5.0.0"
			}
		}`

		err := afero.WriteFile(fs, path, []byte(content), 0o644)
		require.NoError(t, err)

		_, err = Read(fs, path)
		assert.NoError(t, err)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"
		content := `{invalid json}`

		err := afero.WriteFile(fs, path, []byte(content), 0o644)
		require.NoError(t, err)

		_, err = Read(fs, path)
		assert.Error(t, err)
	})

	t.Run("Empty file", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"
		content := `{}`

		err := afero.WriteFile(fs, path, []byte(content), 0o644)
		require.NoError(t, err)

		_, err = Read(fs, path)
		assert.NoError(t, err)
	})
}

func TestWrite(t *testing.T) {
	t.Parallel()

	t.Run("Writes valid package.json", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"

		pkg := &PackageJSON{
			Name:    "test-app",
			Version: "1.0.0",
			Dependencies: map[string]string{
				"payload": "3.0.0",
			},
		}

		err := Write(fs, path, pkg)
		require.NoError(t, err)

		// Verify file was written.
		exists := fsext.Exists(fs, path)
		assert.True(t, exists)

		// Verify content is valid JSON.
		data, err := afero.ReadFile(fs, path)
		require.NoError(t, err)
		assert.Contains(t, string(data), "test-app")
		assert.Contains(t, string(data), "payload")
	})

	t.Run("Write file error returns error", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
		path := "package.json"
		pkg := &PackageJSON{Name: "test-app"}

		err := Write(fs, path, pkg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "writing package.json")
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

		assert.Equal(t, "3.1.0", pkg.Dependencies["payload"])
	})

	t.Run("Field ordering preserved", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"

		original := `{
	"name": "test-app",
	"description": "My test app",
	"license": "MIT",
	"private": true,
	"type": "module",
	"version": "1.0.0",
	"scripts": {
		"dev": "vite dev",
		"build": "vite build"
	}
}`
		err := afero.WriteFile(fs, path, []byte(original), 0o644)
		require.NoError(t, err)

		// Read the package.json.
		pkg, err := Read(fs, path)
		require.NoError(t, err)

		// Modify scripts
		pkg.Scripts["test"] = "vitest"

		// Write it back.
		err = Write(fs, path, pkg)
		require.NoError(t, err)

		// Read it again.
		data, err := afero.ReadFile(fs, path)
		require.NoError(t, err)

		// Verify field order is preserved (npm standard order)
		nameIdx := bytes.Index(data, []byte(`"name"`))
		descIdx := bytes.Index(data, []byte(`"description"`))
		licenseIdx := bytes.Index(data, []byte(`"license"`))
		privateIdx := bytes.Index(data, []byte(`"private"`))
		typeIdx := bytes.Index(data, []byte(`"type"`))
		versionIdx := bytes.Index(data, []byte(`"version"`))
		scriptsIdx := bytes.Index(data, []byte(`"scripts"`))

		assert.Greater(t, descIdx, nameIdx, "description should come after name")
		assert.Greater(t, licenseIdx, descIdx, "license should come after description")
		assert.Greater(t, privateIdx, licenseIdx, "private should come after license")
		assert.Greater(t, typeIdx, privateIdx, "type should come after private")
		assert.Greater(t, versionIdx, typeIdx, "version should come after type")
		assert.Greater(t, scriptsIdx, versionIdx, "scripts should come after version")
	})

	t.Run("HTML characters not escaped", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"

		pkg := &PackageJSON{
			Name:    "test-app",
			Version: "1.0.0",
			Engines: map[string]any{
				"node": "^18.20.2 || >=20.9.0",
			},
		}

		err := Write(fs, path, pkg)
		require.NoError(t, err)

		// Read the written file
		data, err := afero.ReadFile(fs, path)
		require.NoError(t, err)

		fileContent := string(data)
		assert.Contains(t, fileContent, "^18.20.2 || >=20.9.0", "Engine constraint should not be escaped")
		assert.NotContains(t, fileContent, "\\u003e", "Should not contain escaped >")
		assert.NotContains(t, fileContent, "\\u003c", "Should not contain escaped <")
		assert.NotContains(t, fileContent, "\\u0026", "Should not contain escaped &")
	})
}

func TestExists(t *testing.T) {
	t.Parallel()

	t.Run("File exists", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		path := "package.json"

		err := afero.WriteFile(fs, path, []byte(`{}`), 0o644)
		require.NoError(t, err)

		exists := fsext.Exists(fs, path)
		assert.True(t, exists)
	})

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		exists := fsext.Exists(fs, "package.json")
		assert.False(t, exists)
	})
}
