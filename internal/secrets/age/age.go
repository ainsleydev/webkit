package age

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"filippo.io/age"

	"github.com/ainsleydev/webkit/internal/config"
)

const (
	// KeyFileName is the name of the age private key file.
	KeyFileName = "age.key"
	// KeyEnvVar is the environment variable name for the age key.
	KeyEnvVar = "SOPS_AGE_KEY"
)

// ReadIdentity returns the age identity key from environment or file.
//
// Checks in order:
// 1. SOPS_AGE_KEY environment variable.
// 2. ~/.config/webkit/age.key (local dev)
func ReadIdentity() (*age.X25519Identity, error) {
	key := ""
	source := ""

	// Check environment variable first (used in CI/CD)
	if envKey := os.Getenv(KeyEnvVar); envKey != "" {
		key = envKey
		source = fmt.Sprintf("%s environment variable", KeyEnvVar)
	} else {
		// Try reading from config file
		data, err := config.Read(KeyFileName)
		if err != nil {
			path, _ := config.Path(KeyFileName)
			return nil, fmt.Errorf("reading age key from %s: %w", path, err)
		}
		key = string(data)
		keyPath, _ := config.Path(KeyFileName)
		source = keyPath
	}

	if key == "" {
		return nil, errors.New("no SOPS_AGE_KEY key found")
	}

	// Sometimes editors add some random stuff.
	key = strings.ReplaceAll(strings.TrimSpace(key), "\n", "")
	identity, err := age.ParseX25519Identity(strings.TrimSpace(key))
	if err != nil {
		return nil, fmt.Errorf("invalid age key format in %s: %w", source, err)
	}

	return identity, nil
}

// WritePrivateKey writes an age private key to the config directory.
func WritePrivateKey(key string) error {
	// Validate it's a proper age key before writing
	if _, err := age.ParseX25519Identity(key); err != nil {
		return fmt.Errorf("invalid age key format: %w", err)
	}

	if err := config.Write(KeyFileName, []byte(key), os.ModePerm); err != nil {
		return fmt.Errorf("writing age key: %w", err)
	}

	return nil
}

// extractPublicKey extracts the public key from an age private key
func extractPublicKey(privateKey string) (string, error) {
	identity, err := age.ParseX25519Identity(privateKey)
	if err != nil {
		return "", fmt.Errorf("parsing age identity: %w", err)
	}
	return identity.Recipient().String(), nil
}
