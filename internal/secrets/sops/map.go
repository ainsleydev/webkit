package sops

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// DecryptFileToMap decrypts a file using the provided Decrypter and
// returns the content as a map[string]any.
func DecryptFileToMap(ec EncrypterDecrypter, filePath string) (map[string]any, error) {
	if err := ec.Decrypt(filePath); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sops file: %w", err)
	}

	var data map[string]any
	if err = yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse sops content: %w", err)
	}

	// Make sure the file is encrypted after we've
	return data, ec.Encrypt(filePath)
}

func DecryptFileToMapStdout(ec EncrypterDecrypter, filePath string) (map[string]any, error) {
	// Use the Client's runSopsCommand or similar to capture stdout
	out, err := ec.(*Client).runSopsCommand("--decrypt", filePath)
	if err != nil {
		// Example: detect SOPS metadata not found
		if out != "" && strings.Contains(out, "sops metadata not found") {
			return nil, ErrNotEncrypted
		}
		return nil, fmt.Errorf("sops decrypt failed: %s: %w", out, err)
	}

	var data map[string]any
	if err := yaml.Unmarshal([]byte(out), &data); err != nil {
		return nil, fmt.Errorf("failed to parse sops content: %w", err)
	}

	return data, nil
}
