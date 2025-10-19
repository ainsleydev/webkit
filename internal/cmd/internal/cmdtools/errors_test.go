package cmdtools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitError_Error(t *testing.T) {
	t.Parallel()

	err := ExitWithCode(1)
	assert.Equal(t, err.Code, 1)
	assert.Equal(t, err.Error(), "")
}
