package enforce

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/printer"
)

func runFailure(f func()) (string, int) {
	var buf bytes.Buffer
	var code int = -1

	oldConsole := console
	console = printer.New(&buf)
	defer func() { console = oldConsole }()

	oldExit := exit
	exit = func(c int) { code = c }
	defer func() { exit = oldExit }()

	f()

	return buf.String(), code
}

func runSuccess(f func()) string {
	var buf bytes.Buffer

	oldConsole := console
	console = printer.New(&buf)
	defer func() { console = oldConsole }()

	oldExit := exit
	exit = func(c int) { panic("exit should not be called") }
	defer func() { exit = oldExit }()

	f()
	return buf.String()
}

func TestEqual(t *testing.T) {
	t.Run("Passing", func(t *testing.T) {
		out := runSuccess(func() {
			Equal(2, 2, "should pass")
		})
		assert.Empty(t, out)
	})

	t.Run("Failing", func(t *testing.T) {
		out, code := runFailure(func() {
			Equal(1, 2, "values mismatch")
		})
		assert.Contains(t, out, "values mismatch")
		assert.Contains(t, out, "Got:  1")
		assert.Contains(t, out, "Want: 2")
		assert.Equal(t, 1, code)
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("Passing", func(t *testing.T) {
		out := runSuccess(func() {
			NotEqual(1, 2, "should pass")
		})
		assert.Empty(t, out)
	})

	t.Run("Failing", func(t *testing.T) {
		out, code := runFailure(func() {
			NotEqual(1, 1, "values must not match")
		})
		assert.Contains(t, out, "values must not match")
		assert.Contains(t, out, "Got:  1")
		assert.Contains(t, out, "must not equal 1")
		assert.Equal(t, 1, code)
	})
}

func TestTrue(t *testing.T) {
	t.Run("Passing", func(t *testing.T) {
		out := runSuccess(func() {
			True(true, "should pass")
		})
		assert.Empty(t, out)
	})

	t.Run("Failing", func(t *testing.T) {
		out, code := runFailure(func() {
			True(false, "condition false")
		})
		assert.Contains(t, out, "condition false")
		assert.Equal(t, 1, code)
	})
}

func TestNoError(t *testing.T) {
	t.Run("Passing", func(t *testing.T) {
		out := runSuccess(func() {
			NoError(nil, "should pass")
		})
		assert.Empty(t, out)
	})

	t.Run("Failing", func(t *testing.T) {
		out, code := runFailure(func() {
			NoError(errors.New("oops"), "error occurred")
		})
		assert.Contains(t, out, "error occurred")
		assert.Contains(t, out, "Error: oops")
		assert.Equal(t, 1, code)
	})
}

func TestNotNil(t *testing.T) {
	t.Run("Passing", func(t *testing.T) {
		out := runSuccess(func() {
			var x = 123
			NotNil(&x, "should not fail")
		})
		assert.Empty(t, out)
	})

	t.Run("Failing", func(t *testing.T) {
		out, code := runFailure(func() {
			var x *int = nil
			NotNil(x, "value must not be nil")
		})
		assert.Contains(t, out, "value must not be nil")
		assert.Equal(t, 1, code)
	})
}
