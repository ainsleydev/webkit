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
		t.Setenv("BACK_BLAZE_BUCKET", "bucket")
		t.Setenv("BACK_BLAZE_KEY_ID", "id")
		t.Setenv("BACK_BLAZE_APPLICATION_KEY", "appkey")
		t.Setenv("GITHUB_TOKEN", "token")

		cfg, err := ParseTFEnvironment()
		assert.NoError(t, err)
		assert.Equal(t, "key", cfg.DigitalOceanAPIKey)
		assert.Equal(t, "access", cfg.DigitalOceanSpacesAccessKey)
		assert.Equal(t, "secret", cfg.DigitalOceanSpacesSecretKey)
		assert.Equal(t, "bucket", cfg.BackBlazeBucket)
		assert.Equal(t, "id", cfg.BackBlazeKeyID)
		assert.Equal(t, "appkey", cfg.BackBlazeApplicationKey)
		assert.Equal(t, "token", cfg.GithubToken)
	})

	t.Run("Failure", func(t *testing.T) {
		teardownEnv(t) // Sanity check
		_, err := ParseTFEnvironment()
		assert.Error(t, err)
	})
}
