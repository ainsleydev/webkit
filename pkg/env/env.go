package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

// Environment definitions.
const (
	// Development env definition.
	Development string = "development"
	// Staging env definition.
	Staging = "staging"
	// Production env definition.
	Production = "production"
)

// Common keys
const (
	// AppEnvironmentKey is the key for the app environment, i.e. prod/dev
	AppEnvironmentKey = "APP_ENV"
)

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

// Get provides a default override for environment vars as a second argument.
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetOrError will attempt to get the value from the environment variable and
// return an error if not available.
func GetOrError(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("%s is empty", key)
}

// AppEnvironment returns the app environment.
func AppEnvironment() string {
	return strings.ToLower(Get(AppEnvironmentKey, Development))
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
