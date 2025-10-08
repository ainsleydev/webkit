package sops

import (
	"fmt"

	"go.mozilla.org/sops/decrypt"
)

// IsContentEncrypted checks if file contents is SOPS encrypted.
func IsContentEncrypted(content []byte) bool {
	// Attempt to decrypt the content
	_, err := decrypt.Data(content, "json")
	fmt.Print(err)
	// If decryption is successful, the file is encrypted
	return err == nil

}

// DecryptValue decrypts a SOPS file and extracts a specific key.
// This is the main function used during `webkit update` to resolve
// environment variables with source type "sops".
//
// Example:
//
//	value, err := sops.DecryptValue("secrets/production.yaml", "PAYLOAD_SECRET")
//
// DecryptValue decrypts a SOPS file and extracts a specific key.
// This is the main function used during `webkit update` to resolve
// environment variables with source type "sops".
//
// Example:
//
//	value, err := sops.DecryptValue("secrets/production.yaml", "PAYLOAD_SECRET")
//func DecryptValue(filePath, key string) (string, error) {
//	// Get age key and validate it
//	ageKey, err := AgeKeyRead()
//	if err != nil {
//		return "", err
//	}
//
//	// Validate it's a proper age key
//	if _, err := age.ParseX25519Identity(ageKey); err != nil {
//		keyPath, _ := config.Path(AgeKeyFileName)
//		return "", fmt.Errorf("invalid age key format in %s: %w", keyPath, err)
//	}
//
//	// Set in environment for SOPS
//	os.Setenv(AgeKeyEnvVar, ageKey)
//	defer os.Unsetenv(AgeKeyEnvVar)
//
//	// Decrypt the entire file using SOPS
//	decrypted, err := decrypt.File(filePath, "yaml")
//	if err != nil {
//		return "", fmt.Errorf("decrypting %s: %w", filePath, err)
//	}
//
//	// Parse YAML and extract the key
//	var data map[string]interface{}
//	if err := yaml.Unmarshal(decrypted, &data); err != nil {
//		return "", fmt.Errorf("parsing decrypted YAML from %s: %w", filePath, err)
//	}
//
//	// Extract value
//	value, ok := data[key]
//	if !ok {
//		return "", fmt.Errorf("key %q not found in %s", key, filePath)
//	}
//
//	// Convert to string
//	strValue, ok := value.(string)
//	if !ok {
//		return "", fmt.Errorf("key %q in %s is not a string value", key, filePath)
//	}
//
//	return strValue, nil
//}
