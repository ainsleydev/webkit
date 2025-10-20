package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	t.Parallel()

	got := Info()

	t.Log("Formatted Info")
	{
		assert.Contains(t, got, "WebKit")
		assert.Contains(t, got, "Commit:")
		assert.Contains(t, got, "Built:")
		assert.Contains(t, got, "Built by:")
	}

	t.Log("Has Value")
	{
		assert.NotEmpty(t, Version)
		assert.Contains(t, got, Version)
	}

	t.Log("Multiline String")
	{
		lines := strings.Split(got, "\n")
		assert.Len(t, lines, 4)
	}
}
