package secrets

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

var ScaffoldCmd = &cli.Command{
	Name:   "scaffold",
	Usage:  "Generate empty SOPS secret files and configuration",
	Action: cmdtools.Wrap(Scaffold),
}

// Scaffold generates the basic SOPS secret file structure.
// This creates empty secret files and SOPS configuration without
// parsing app.json.
func Scaffold(_ context.Context, input cmdtools.CommandInput) error {
	if err := generateSOPSConfig(input.Generator()); err != nil {
		return errors.Wrap(err, "generating sops config")
	}

	for _, enviro := range env.All {
		path := filepath.Join("resources", "secrets", fmt.Sprintf("%s.yaml", enviro))

		// If we generate a file that has YAML commentary in the file,
		// SOPS will encrypt the comments when Encrypt() is called,
		// malforming the file.
		err := input.
			Generator().
			Bytes(path, make([]byte, 0), scaffold.WithScaffoldMode(), scaffold.WithoutNotice())
		if err != nil {
			return fmt.Errorf("generating %s: %w", path, err)
		}
	}

	return nil
}

// generateSOPSConfig creates the .sops.yaml configuration file which tells
// sops how to encrypt and decrypt files without specifying rules or keys
// everytime we call the cmd.
func generateSOPSConfig(gen scaffold.Generator) error {
	config := map[string]any{
		"creation_rules": []map[string]any{
			{
				"path_regex": `secrets/.*\.yaml$`,
				"age":        secrets.AgePublicKey,
			},
		},
	}
	return gen.YAML(filepath.Join("resources", ".sops.yaml"), config)
}
