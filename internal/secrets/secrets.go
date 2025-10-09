package secrets

import "path/filepath"

// AgePublicKey is the public key for encrypting SOPS files.
const AgePublicKey = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"

// FilePath defines the path where SOPS encrypted YAML files
// reside in the Webkit app. Needs a base path prepended.
var FilePath = filepath.Join("resources", "secrets")
