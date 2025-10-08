package operations

import (
	"context"
	"fmt"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/pkg/env"
)

// CreateSecretFiles generates the basic SOPS secret file structure.
// This creates empty secret files and SOPS configuration without
// parsing app.json.
func CreateSecretFiles(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(afero.NewBasePathFs(input.FS, "resources"))

	// Generate .sops.yaml configuration.
	if err := generateSOPSConfig(gen); err != nil {
		return fmt.Errorf("generating .sops.yaml: %w", err)
	}

	// Generate empty secret files for each environment.
	environments := []string{env.Development, env.Staging, env.Production}
	for _, enviro := range environments {
		filePath := fmt.Sprintf("secrets/%s.yaml", enviro)
		if err := generateEmptySecretFile(gen, filePath, enviro); err != nil {
			return fmt.Errorf("generating %s: %w", filePath, err)
		}
	}

	return nil
}

// generateSOPSConfig creates the .sops.yaml configuration file
func generateSOPSConfig(gen scaffold.Generator) error {
	config := map[string]interface{}{
		"creation_rules": []map[string]interface{}{
			{
				"path_regex": `secrets/.*\.yaml$`,
				"age":        "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			},
		},
	}

	return gen.YAML(".sops.yaml", config, scaffold.WithScaffoldMode())
}

// generateEmptySecretFile creates an empty secret file with instructions
func generateEmptySecretFile(gen scaffold.Generator, path string, env string) error {
	content := fmt.Sprintf(`# %s environment secrets
#
# Add your secret values below in the format:
#   KEY_NAME: "value"
#
# After adding secrets, encrypt this file:
#   sops --encrypt --in-place %s
#
# To edit encrypted secrets:
#   sops %s
#
# To decrypt (not recommended, use 'sops' to edit):
#   sops --decrypt %s
#
# WebKit will automatically decrypt these during 'webkit update'

`, env, path, path, path)

	return gen.Bytes(path, []byte(content), scaffold.WithScaffoldMode())
}
