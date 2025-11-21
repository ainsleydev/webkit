package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecCmd(t *testing.T) {
	t.Parallel()

	t.Run("Command Registered", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, ExecCmd)
		assert.Equal(t, "exec", ExecCmd.Name)
		assert.Equal(t, "Execute arbitrary Terraform commands", ExecCmd.Usage)
		assert.Contains(t, ExecCmd.Description, "Examples:")
		assert.Contains(t, ExecCmd.Description, "state list")
		assert.NotNil(t, ExecCmd.Action)
	})
}
