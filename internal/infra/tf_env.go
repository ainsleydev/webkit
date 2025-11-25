package infra

import (
	"log/slog"

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
	HetznerToken                string `env:"HETZNER_TOKEN,required"`
	BackBlazeBucket             string `env:"BACK_BLAZE_BUCKET,required"`
	BackBlazeKeyID              string `env:"BACK_BLAZE_KEY_ID,required"`
	BackBlazeApplicationKey     string `env:"BACK_BLAZE_APPLICATION_KEY,required"`
	TursoToken                  string `env:"TURSO_TOKEN,required"`
	GithubToken                 string `env:"GITHUB_TOKEN,required"`
	GithubTokenClassic          string `env:"GITHUB_TOKEN_CLASSIC,required"`
	SlackBotToken               string `env:"SLACK_BOT_TOKEN,required"`
	SlackUserToken              string `env:"SLACK_USER_TOKEN,required"`
	SlackWebhookURL             string `env:"SLACK_WEBHOOK_URL"`
	PeekapingEndpoint           string `env:"PEEKAPING_ENDPOINT"`
	PeekapingAPIKey             string `env:"PEEKAPING_API_KEY"`
}

// ParseTFEnvironment reads and validates Terraform-related
// environment variables.
func ParseTFEnvironment() (TFEnvironment, error) {
	cfg, err := env.ParseAs[TFEnvironment]()
	if err != nil {
		return TFEnvironment{}, errors.Wrap(err, "parsing terraform environment")
	}

	// Log which credentials are present (not the values, just presence)
	slog.Info("Parsed Terraform environment variables",
		slog.Bool("has_do_token", cfg.DigitalOceanAPIKey != ""),
		slog.Bool("has_do_spaces_creds", cfg.DigitalOceanSpacesAccessKey != "" && cfg.DigitalOceanSpacesSecretKey != ""),
		slog.Bool("has_hetzner_token", cfg.HetznerToken != ""),
		slog.Bool("has_b2_creds", cfg.BackBlazeBucket != "" && cfg.BackBlazeKeyID != "" && cfg.BackBlazeApplicationKey != ""),
		slog.Bool("has_turso_token", cfg.TursoToken != ""),
		slog.Bool("has_github_tokens", cfg.GithubToken != "" && cfg.GithubTokenClassic != ""),
		slog.Bool("has_slack_tokens", cfg.SlackBotToken != "" && cfg.SlackUserToken != ""),
		slog.Bool("has_slack_webhook", cfg.SlackWebhookURL != ""),
		slog.Bool("has_peekaping_endpoint", cfg.PeekapingEndpoint != ""),
		slog.Bool("has_peekaping_api_key", cfg.PeekapingAPIKey != ""),
	)

	return cfg, nil
}

// varStrings maps the environment to Terraform variable strings
// to pass to the execer.
func (t *TFEnvironment) varStrings() []string {
	vars := []string{
		"do_token=" + t.DigitalOceanAPIKey,
		"do_spaces_access_id=" + t.DigitalOceanSpacesAccessKey,
		"do_spaces_secret_key=" + t.DigitalOceanSpacesSecretKey,
		"hetzner_token=" + t.HetznerToken,
		"b2_application_key=" + t.BackBlazeApplicationKey,
		"b2_application_key_id=" + t.BackBlazeKeyID,
		"turso_api_token=" + t.TursoToken,
		"github_token=" + t.GithubToken,
		"github_token_classic=" + t.GithubTokenClassic,
		"slack_bot_token=" + t.SlackBotToken,
		"slack_user_token=" + t.SlackUserToken,
	}

	// Only include Peekaping credentials if they are configured.
	// This prevents provider initialization when monitoring is not in use.
	if t.PeekapingEndpoint != "" {
		vars = append(vars, "peekaping_endpoint="+t.PeekapingEndpoint)
	}
	if t.PeekapingAPIKey != "" {
		vars = append(vars, "peekaping_api_key="+t.PeekapingAPIKey)
	}

	return vars
}
