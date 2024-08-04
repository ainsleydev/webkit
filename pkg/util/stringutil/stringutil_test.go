package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()
		assert.False(t, IsNotEmpty(nil))
		assert.False(t, IsNotEmpty(ptr.StringPtr("")))
	})

	t.Run("Non Empty", func(t *testing.T) {
		t.Parallel()
		assert.True(t, IsNotEmpty(ptr.StringPtr("hello")))
	})
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()
		assert.True(t, IsEmpty(nil))
		assert.True(t, IsEmpty(ptr.StringPtr("")))
	})

	t.Run("Non Empty", func(t *testing.T) {
		t.Parallel()
		assert.False(t, IsEmpty(ptr.StringPtr("hello")))
	})
}

func TestRemoveDuplicateWhitespace(t *testing.T) {
	t.Parallel()
	in := `   <source srcset="https://example.com/image.jpg"
		type="image/webp"

	/>  `
	want := `<source srcset="https://example.com/image.jpg" type="image/webp" />`
	got := RemoveDuplicateWhitespace(in)
	assert.Equal(t, want, got)
}

func TestFormatHTML(t *testing.T) {
	t.Parallel()
	in := `<picture > <img src=\"https://example.com/image.jpg\" alt=\"Alternative\" /> </picture>`
	want := "<picture>\n  <img src=\\\"https://example.com/image.jpg\\\" alt=\\\"Alternative\\\" />\n</picture>"
	got := FormatHTML(in)
	assert.Equal(t, want, got)
}
