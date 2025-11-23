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
			"monitoring": {
				"peekaping_endpoint": "https://peekaping.example.com",
				"http_monitors": {
					"HTTP - example.com": {"id": "abc123", "name": "HTTP - example.com"}
				},
				"dns_monitors": {},
				"push_monitors": {},
				"status_page_url": "https://peekaping.example.com/status/test"
			},
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
		assert.Equal(t, "https://peekaping.example.com", got.Monitoring.PeekapingEndpoint)
		assert.Equal(t, "https://peekaping.example.com/status/test", got.Monitoring.StatusPageURL)
		assert.Equal(t, "alerts-test", got.Slack.ChannelName)
		assert.Equal(t, "C123456", got.Slack.ChannelID)
		assert.Contains(t, got.Monitoring.HTTPMonitors, "HTTP - example.com")
		monitor := got.Monitoring.HTTPMonitors["HTTP - example.com"].(map[string]any)
		assert.Equal(t, "abc123", monitor["id"])
	})
}
