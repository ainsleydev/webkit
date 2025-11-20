//go:build !race

package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTFEnvironment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		defer teardownEnv(t)

		t.Setenv("DO_API_KEY", "key")
		t.Setenv("DO_SPACES_ACCESS_KEY", "access")
		t.Setenv("DO_SPACES_SECRET_KEY", "secret")
		t.Setenv("HETZNER_TOKEN", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
		t.Setenv("BACK_BLAZE_BUCKET", "bucket")
		t.Setenv("BACK_BLAZE_KEY_ID", "id")
		t.Setenv("BACK_BLAZE_APPLICATION_KEY", "appkey")
		t.Setenv("TURSO_TOKEN", "turso-test-token")
		t.Setenv("GITHUB_TOKEN", "token")
		t.Setenv("GITHUB_TOKEN_CLASSIC", "token")
		t.Setenv("SLACK_BOT_TOKEN", "xoxb-test-token")
		t.Setenv("SLACK_USER_TOKEN", "xoxp-test-token")
		t.Setenv("SLACK_WEBHOOK_URL", "https://hooks.slack.com/services/test")
		t.Setenv("PEEKAPING_ENDPOINT", "https://uptime.test.dev")
		t.Setenv("PEEKAPING_EMAIL", "test@example.com")
		t.Setenv("PEEKAPING_PASSWORD", "test-password")

		cfg, err := ParseTFEnvironment()
		assert.NoError(t, err)
		assert.Equal(t, "key", cfg.DigitalOceanAPIKey)
		assert.Equal(t, "access", cfg.DigitalOceanSpacesAccessKey)
		assert.Equal(t, "secret", cfg.DigitalOceanSpacesSecretKey)
		assert.Equal(t, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", cfg.HetznerToken)
		assert.Equal(t, "bucket", cfg.BackBlazeBucket)
		assert.Equal(t, "id", cfg.BackBlazeKeyID)
		assert.Equal(t, "appkey", cfg.BackBlazeApplicationKey)
		assert.Equal(t, "turso-test-token", cfg.TursoToken)
		assert.Equal(t, "token", cfg.GithubToken)
		assert.Equal(t, "token", cfg.GithubTokenClassic)
		assert.Equal(t, "xoxb-test-token", cfg.SlackBotToken)
		assert.Equal(t, "xoxp-test-token", cfg.SlackUserToken)
		assert.Equal(t, "https://hooks.slack.com/services/test", cfg.SlackWebhookURL)
		assert.Equal(t, "https://uptime.test.dev", cfg.PeekapingEndpoint)
		assert.Equal(t, "test@example.com", cfg.PeekapingEmail)
		assert.Equal(t, "test-password", cfg.PeekapingPassword)
	})

	t.Run("Failure", func(t *testing.T) {
		teardownEnv(t) // Sanity check
		_, err := ParseTFEnvironment()
		assert.Error(t, err)
	})
}
