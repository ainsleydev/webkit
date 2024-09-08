package storage

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageProviderPersistence(t *testing.T) {
	t.Parallel()

	providers := map[string]Provider{
		"OS": setupOSStorage(t),
		"S3": setupS3StorageForPersistenceTest(t),
	}

	for name, provider := range providers {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testPersistence(t, provider)
		})
	}
}

func testPersistence(t *testing.T, provider Provider) {
	t.Helper()

	testFile := "test-file.txt"
	testContent := []byte("Hello, Storage!")
	ctx := context.Background()

	// Test Upload
	err := provider.Upload(ctx, testFile, bytes.NewReader(testContent))
	require.NoError(t, err)

	// Test Exists
	exists, err := provider.Exists(ctx, testFile)
	require.NoError(t, err)
	assert.True(t, exists)

	// Test List
	files, err := provider.List(ctx, "")
	require.NoError(t, err)
	assert.Contains(t, files, testFile)

	// Test Stat
	info, err := provider.Stat(ctx, testFile)
	require.NoError(t, err)
	assert.Equal(t, int64(len(testContent)), info.Size)
	assert.False(t, info.IsDir)
	assert.WithinDuration(t, time.Now(), info.LastModified, 2*time.Second)

	// Test Download
	reader, err := provider.Download(ctx, testFile)
	require.NoError(t, err)
	downloadedContent, err := io.ReadAll(reader)
	require.NoError(t, err)
	reader.Close()
	assert.Equal(t, testContent, downloadedContent)

	// Test Delete
	err = provider.Delete(ctx, testFile)
	require.NoError(t, err)

	// Verify file no longer exists
	exists, err = provider.Exists(ctx, testFile)
	require.NoError(t, err)
	assert.False(t, exists)

	// Verify file is not in the list
	files, err = provider.List(ctx, "")
	require.NoError(t, err)
	assert.NotContains(t, files, testFile)

	// Attempt to download non-existent file
	_, err = provider.Download(ctx, testFile)
	assert.Error(t, err)

	// Attempt to stat non-existent file
	_, err = provider.Stat(ctx, testFile)
	assert.Error(t, err)
}
