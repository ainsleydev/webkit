package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("Stdout mode", func(t *testing.T) {
		var buf bytes.Buffer
		p := printer.New(&buf)

		err := run(context.Background(), p, "", true)

		assert.NoError(t, err)
	})

	t.Run("File output mode", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "test-schema.json")

		var buf bytes.Buffer
		p := printer.New(&buf)

		err := run(context.Background(), p, outputPath, false)

		require.NoError(t, err)

		data, readErr := os.ReadFile(outputPath)
		require.NoError(t, readErr)
		assert.NotEmpty(t, data)
	})

	t.Run("Creates directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "nested", "dir", "schema.json")

		var buf bytes.Buffer
		p := printer.New(&buf)

		err := run(context.Background(), p, outputPath, false)

		require.NoError(t, err)

		dir := filepath.Dir(outputPath)
		_, statErr := os.Stat(dir)
		require.NoError(t, statErr)

		data, readErr := os.ReadFile(outputPath)
		require.NoError(t, readErr)
		assert.NotEmpty(t, data)
	})
}
