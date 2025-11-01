package cmdtools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitError_Error(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		code int
	}{
		"Code 0": {code: 0},
		"Code 1": {code: 1},
		"Code 2": {code: 2},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := ExitWithCode(test.code)
			assert.Equal(t, test.code, err.Code)
			assert.Equal(t, "", err.Error())
		})
	}
}

func TestExitWithCode(t *testing.T) {
	t.Parallel()

	err := ExitWithCode(1)
	assert.Equal(t, 1, err.Code)

	// Verify it can be used as an error.
	var _ error = err
}
