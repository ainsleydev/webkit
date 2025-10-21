package secrets

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
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
				//if !slices.Contains(env.All, s) {
				//	return fmt.Errorf("invalid environment: %s", s)
				//}
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "key",
			Usage:   "The key/name of the secret to retrieve",
			Aliases: []string{"k"},
		},
		&cli.BoolFlag{
			Name:    "all",
			Usage:   "Print all secrets from the environment",
			Aliases: []string{"a"},
		},
	},
	Action: cmdtools.Wrap(Get),
}

// Get retrieves a singular or all decrypted secret(s) by environment.
// Use with caution, it displays confidential secrets.
func Get(_ context.Context, input cmdtools.CommandInput) error {
	cmd := input.Command
	enviro := cmd.String("env")
	key := cmd.String("key")
	showAll := cmd.Bool("all")
	client := input.SOPSClient()

	path := filepath.Join(input.BaseDir, secrets.FilePathFromEnv(env.Environment(enviro)))
	vals, err := sops.DecryptFileToMap(client, path)
	if err != nil {
		return errors.Wrap(err, "decoding sops to map")
	}

	switch {
	case showAll:
		input.Printer().Success(fmt.Sprintf("All secrets retrieved for environment: %s\n", enviro))
		var items []string
		for k, v := range vals {
			items = append(items, fmt.Sprintf("%s: %s", k, v))
		}
		input.Printer().List(items)
	case key != "":
		value, ok := vals[key]
		if !ok {
			return fmt.Errorf("key %v not found for env: %s", key, enviro)
		}
		input.Printer().Success(fmt.Sprintf("Successfully retrieved secret from environment: %s\n", enviro))
		input.Printer().Printf("%s=%v", key, value)
	default:
		return fmt.Errorf("either --key or --all must be provided")
	}

	return nil
}
