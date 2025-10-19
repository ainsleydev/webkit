package secrets

import (
	"path/filepath"

	"github.com/ainsleydev/webkit/pkg/env"
)

// AgePublicKey is the public key for encrypting SOPS files.
const AgePublicKey = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"

// FilePath defines the path where SOPS encrypted YAML files
// reside in the Webkit app. Needs a base path prepended.
var FilePath = filepath.Join("resources", "secrets")

// FilePathFromEnv returns a filepath based off the environment.
//
// For example, resources/secrets/{production}.yaml
func FilePathFromEnv(e env.Environment) string {
	return filepath.Join(FilePath, e.String()+".yaml")
}
