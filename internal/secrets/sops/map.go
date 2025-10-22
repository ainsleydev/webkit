package sops

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DecryptFileToMap decrypts a file using the provided Decrypter and
// returns the content as a map[string]any.
func DecryptFileToMap(ec EncrypterDecrypter, filePath string) (map[string]any, error) {
	decryptErr := ec.Decrypt(filePath)
	if decryptErr != nil && !errors.Is(decryptErr, ErrNotEncrypted) {
		return nil, decryptErr
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sops file: %w", err)
	}

	var data map[string]any
	if err = yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse sops content: %w", err)
	}

	// Only re-encrypt if we actually decrypted it.
	if decryptErr == nil {
		if encErr := ec.Encrypt(filePath); encErr != nil {
			return nil, encErr
		}
	}

	return data, nil
}
