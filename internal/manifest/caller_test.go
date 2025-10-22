package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaller(t *testing.T) {
	t.Parallel()

	got := testHelper(t)
	assert.Contains(t, []string{"manifest", "manifest_test"}, got.Package)
	assert.Equal(t, "TestCaller", got.Function)
	assert.NotEmpty(t, got.File)
	assert.Greater(t, got.Line, 0)

	t.Run("UnexportedFunctionIsSkipped", func(t *testing.T) {
		t.Parallel()
		got := unexportedHelper(t)
		assert.Equal(t, "unknown", got.Function)
		assert.Equal(t, "unknown", got.Package)
	})

	t.Run("UnknownWhenNoFrames", func(t *testing.T) {
		t.Parallel()
		got := noFrameCaller()
		assert.Equal(t, "unknown", got.Function)
		assert.Equal(t, "unknown", got.Package)
	})
}

func testHelper(t *testing.T) CallerInfo {
	t.Helper()
	return Caller()
}

func unexportedHelper(t *testing.T) CallerInfo {
	t.Helper()
	return lowerCaseFunction()
}

func lowerCaseFunction() CallerInfo {
	return Caller()
}

func noFrameCaller() CallerInfo {
	return CallerInfo{
		Package:  "unknown",
		Function: "unknown",
	}
}
