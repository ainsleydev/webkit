package sops

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DecryptFileToMap decrypts a file using the provided Decrypter and
// returns the content as a map[string]any.
func DecryptFileToMap(d Decrypter, filePath string) (map[string]any, error) {
	if err := d.Decrypt(filePath); err != nil {
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

	return data, nil
}
