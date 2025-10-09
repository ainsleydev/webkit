package secrets

import "path/filepath"

// FilePath defines the path where SOPS encrypted YAML files
// reside in the Webkit app. Needs a base path prepended.
var FilePath = filepath.Join("resources", "secrets")
