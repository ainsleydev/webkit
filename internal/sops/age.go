package sops

import (
	"fmt"
	"os"

	"filippo.io/age"

	"github.com/ainsleydev/webkit/internal/config"
)

const (
	// AgeKeyFileName is the name of the age private key file
	AgeKeyFileName = "age.key"

	// AgeKeyEnvVar is the environment variable name for the age key
	AgeKeyEnvVar = "SOPS_AGE_KEY"
)

// AgeKeyRead returns the age private key from environment or file.
// Checks in order:
//
// 1. SOPS_AGE_KEY environment variable.
// 2. ~/.config/webkit/age.key (local dev)
func AgeKeyRead() (string, error) {
	key := ""
	source := ""

	// Check environment variable first (used in CI/CD)
	if envKey := os.Getenv(AgeKeyEnvVar); envKey != "" {
		key = envKey
		source = fmt.Sprintf("%s environment variable", AgeKeyEnvVar)
	} else {
		// Try reading from config file
		data, err := config.Read(AgeKeyFileName)
		if err != nil {
			path, _ := config.Path(AgeKeyFileName)
			return "", fmt.Errorf("reading age key from %s: %w", path, err)
		}
		key = string(data)
		keyPath, _ := config.Path(AgeKeyFileName)
		source = keyPath
	}

	// Validate the key (only once, regardless of source)
	if _, err := age.ParseX25519Identity(key); err != nil {
		return "", fmt.Errorf("invalid age key format in %s: %w", source, err)
	}

	return key, nil
}

// AgeKeyWrite writes an age private key to the config directory.
func AgeKeyWrite(key string) error {
	// Validate it's a proper age key before writing
	if _, err := age.ParseX25519Identity(key); err != nil {
		return fmt.Errorf("invalid age key format: %w", err)
	}

	if err := config.Write(AgeKeyFileName, []byte(key), os.ModePerm); err != nil {
		return fmt.Errorf("writing age key: %w", err)
	}

	return nil
}
