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
	BackBlazeKeyID              string `env:"BACK_BLAZE_KEY_ID,required"`
	BackBlazeApplicationKey     string `env:"BACK_BLAZE_APPLICATION_KEY,required"`
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
