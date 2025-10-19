package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashContent(t *testing.T) {
	t.Parallel()

	t.Run("Empty Content", func(t *testing.T) {
		t.Parallel()

		got := HashContent([]byte{})
		assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", got)
	})

	t.Run("Consistent Hashing", func(t *testing.T) {
		t.Parallel()

		data := []byte("test content")
		hash1 := HashContent(data)
		hash2 := HashContent(data)

		assert.Equal(t, hash1, hash2)
	})

	t.Run("Different Content Different Hash", func(t *testing.T) {
		t.Parallel()

		hash1 := HashContent([]byte("content1"))
		hash2 := HashContent([]byte("content2"))
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Known Hash Value", func(t *testing.T) {
		t.Parallel()

		data := []byte("hello world")
		got := HashContent(data)
		want := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
		assert.Equal(t, want, got)
	})
}
