package secrets

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
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
	gen := scaffold.New(afero.NewBasePathFs(input.FS, "resources"), input.Manifest)

	if err := generateSOPSConfig(gen); err != nil {
		return fmt.Errorf("generating .sops.yaml: %w", err)
	}

	for _, enviro := range env.All {
		filePath := fmt.Sprintf("secrets/%s.yaml", enviro)

		// If we generate a file that has YAML commentary in the file,
		// SOPS will encrypt the comments when Encrypt() is called,
		// malforming the file.
		err := gen.Bytes(filePath, make([]byte, 0), scaffold.WithScaffoldMode())
		if err != nil {
			return fmt.Errorf("generating %s: %w", filePath, err)
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
	return gen.YAML(".sops.yaml", config)
}
