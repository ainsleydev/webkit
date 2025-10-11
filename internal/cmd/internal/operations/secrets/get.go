package secrets

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
)

var GetCmd = &cli.Command{
	Name:        "get",
	Usage:       "Retrieve a secret from a specific environment",
	Description: "Fetches the value of a secret from the chosen environment (development, staging, production)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "env",
			Usage:    "Environment to fetch the secret from (development, staging, production)",
			Aliases:  []string{"e"},
			Required: true,
			Validator: func(s string) error {
				if !slices.Contains(env.All, s) {
					return fmt.Errorf("invalid environment: %s", s)
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:     "key",
			Usage:    "The key/name of the secret to retrieve",
			Aliases:  []string{"k"},
			Required: true,
		},
	},
	Action: cmdtools.Wrap(Get),
}

// Get retrieves a singular decrypted secret by environment.
// Use with caution.
func Get(_ context.Context, input cmdtools.CommandInput) error {
	cmd := input.Command
	enviro := cmd.String("env")
	key := cmd.String("key")

	client, err := input.SOPSClient()
	if err != nil {
		return err
	}

	path := filepath.Join(input.BaseDir, secrets.FilePath, enviro+".yaml")
	vals, err := sops.DecryptFileToMap(client, path)
	fmt.Println(vals, err)
	if err != nil {
		return errors.Wrap(err, "decoding sops to map")
	}

	value, ok := vals[key]
	if !ok {
		return fmt.Errorf("key %v not found for env: %s", key, enviro)
	}

	fmt.Println(value)

	return nil
}
