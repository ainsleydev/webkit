package infra

import (
	"github.com/caarlos0/env/v11"
	"github.com/pkg/errors"
)

// TFEnvironment holds the required environment variables for Terraform
// operations. Plan and Apply cannot be ran without them as they
// are backend configs.
type TFEnvironment struct {
	DigitalOceanAPIKey          string `env:"DO_API_KEY,required"`
	DigitalOceanSpacesAccessKey string `env:"DO_SPACES_ACCESS_KEY,required"`
	DigitalOceanSpacesSecretKey string `env:"DO_SPACES_SECRET_KEY,required"`
	BackBlazeBucket             string `env:"BACK_BLAZE_BUCKET,required"`
	BackBlazeKeyID              string `env:"BACK_BLAZE_KEY_ID,required"`
	BackBlazeApplicationKey     string `env:"BACK_BLAZE_APPLICATION_KEY,required"`
	TursoToken                  string `env:"TURSO_TOKEN"`
	GithubToken                 string `env:"GITHUB_TOKEN,required"`
	GithubTokenClassic          string `env:"GITHUB_TOKEN_CLASSIC,required"`
	SlackBotToken               string `env:"SLACK_BOT_TOKEN,required"`
	SlackUserToken              string `env:"SLACK_USER_TOKEN,required"`
}

// ParseTFEnvironment reads and validates Terraform-related
// environment variables.
func ParseTFEnvironment() (TFEnvironment, error) {
	cfg, err := env.ParseAs[TFEnvironment]()
	if err != nil {
		return TFEnvironment{}, errors.Wrap(err, "parsing terraform environment")
	}
	return cfg, nil
}

// varStrings maps the environment to Terraform variable strings
// to pass to the execer.
func (t *TFEnvironment) varStrings() []string {
	return []string{
		"do_token=" + t.DigitalOceanAPIKey,
		"do_spaces_access_id=" + t.DigitalOceanSpacesAccessKey,
		"do_spaces_secret_key=" + t.DigitalOceanSpacesSecretKey,
		"b2_application_key=" + t.BackBlazeApplicationKey,
		"b2_application_key_id=" + t.BackBlazeKeyID,
		"turso_api_token=" + t.TursoToken,
		"github_token=" + t.GithubToken,
		"github_token_classic=" + t.GithubTokenClassic,
		"slack_bot_token=" + t.SlackBotToken,
		"slack_user_token=" + t.SlackUserToken,
	}
}
