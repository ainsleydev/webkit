package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Environment definitions.
const (
	// Development env definition.
	Development Environment = "development"
	// Staging env definition.
	Staging Environment = "staging"
	// Production env definition.
	Production Environment = "production"
)

// Environment represents an application deployment environment.
type Environment string

// String returns the string representation of the environment.
func (e Environment) String() string {
	return string(e)
}

// Short returns a short name for the environment.
func (e Environment) Short() string {
	switch e {
	case Development:
		return "dev"
	case Staging:
		return "staging"
	case Production:
		return "prod"
	default:
		return "unknown"
	}
}

// Common keys
const (
	// AppEnvironmentKey is the key for the app environment, i.e. prod/dev
	AppEnvironmentKey = "APP_ENV"
)

// All defines all environments combined.
var All = []Environment{
	Development,
	Staging,
	Production,
}

// ParseConfig loads the environment variables from the .env file and parses the
// environment variables into the provided struct. It returns an error if the
// .env file cannot be loaded or if the environment variables cannot be parsed.
//
//	For example:
//
//	type Config struct {
//		Home         string         `env:"HOME"`
//		Port         int            `env:"PORT" envDefault:"3000"`
//		Password     string         `env:"PASSWORD,unset"`
//		IsProduction bool           `env:"PRODUCTION"`
//		Duration     time.Duration  `env:"DURATION"`
//		Hosts        []string       `env:"HOSTS" envSeparator:":"`
//		TempFolder   string         `env:"TEMP_FOLDER,expand" envDefault:"${HOME}/tmp"`
//		StringInts   map[string]int `env:"MAP_STRING_INT"`
//	}
func ParseConfig(cfg any, filenames ...string) error {
	if IsDevelopment() {
		if err := godotenv.Load(filenames...); err != nil {
			return err
		}
	}
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}

// Get retrieves an environment variable value or returns the fallback if not set.
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetOrError retrieves an environment variable value or returns an error if not set.
func GetOrError(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("%s is empty", key)
}

// AppEnvironment returns the app environment.
func AppEnvironment() Environment {
	return Environment(strings.ToLower(Get(AppEnvironmentKey, Development.String())))
}

// IsDevelopment returns whether we are running the app in development.
func IsDevelopment() bool {
	e := AppEnvironment()
	return e == Development || e == ""
}

// IsStaging returns whether we are running the app in staging.
func IsStaging() bool {
	return AppEnvironment() == Staging
}

// IsProduction returns whether we are running the app in production.
func IsProduction() bool {
	return AppEnvironment() == Production
}
