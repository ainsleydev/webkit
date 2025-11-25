package outputs

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		got := Load(fs)
		assert.Nil(t, got)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(".webkit", 0o755)
		require.NoError(t, err)
		err = afero.WriteFile(fs, ".webkit/outputs.json", []byte("invalid json"), 0o644)
		require.NoError(t, err)

		got := Load(fs)
		assert.Nil(t, got)
	})

	t.Run("Valid outputs", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		outputsJSON := `{
			"peekaping": {
				"endpoint": "https://peekaping.example.com",
				"project_tag": "abc-123-def-456"
			},
			"monitors": [
				{"id": "abc123", "name": "HTTP - example.com", "type": "http"},
				{"id": "def456", "name": "DNS - example.com", "type": "dns"}
			],
			"slack": {
				"channel_name": "alerts-test",
				"channel_id": "C123456"
			}
		}`

		err := fs.MkdirAll(".webkit", 0o755)
		require.NoError(t, err)
		err = afero.WriteFile(fs, ".webkit/outputs.json", []byte(outputsJSON), 0o644)
		require.NoError(t, err)

		got := Load(fs)
		require.NotNil(t, got)
		assert.Equal(t, "https://peekaping.example.com", got.Peekaping.Endpoint)
		assert.Equal(t, "abc-123-def-456", got.Peekaping.ProjectTag)
		assert.Equal(t, "alerts-test", got.Slack.ChannelName)
		assert.Equal(t, "C123456", got.Slack.ChannelID)
		require.Len(t, got.Monitors, 2)
		assert.Equal(t, "abc123", got.Monitors[0].ID)
		assert.Equal(t, "HTTP - example.com", got.Monitors[0].Name)
		assert.Equal(t, "http", got.Monitors[0].Type)
	})
}
