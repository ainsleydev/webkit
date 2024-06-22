package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Empty", func(t *testing.T) {
		assert.False(t, IsNotEmpty(nil))
		assert.False(t, IsNotEmpty(ptr.StringPtr("")))
	})

	t.Run("Non Empty", func(t *testing.T) {
		assert.True(t, IsNotEmpty(ptr.StringPtr("hello")))
	})
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Empty", func(t *testing.T) {
		assert.True(t, IsEmpty(nil))
		assert.True(t, IsEmpty(ptr.StringPtr("")))
	})

	t.Run("Non Empty", func(t *testing.T) {
		assert.False(t, IsEmpty(ptr.StringPtr("hello")))
	})
}
