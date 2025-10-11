package sops

import (
	"fmt"
	"os"

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

	fmt.Println(string(content))
	var data map[string]any
	if err = yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse sops content: %w", err)
	}

	// Make sure the file is encrypted after we've
	return data, ec.Encrypt(filePath)
}
