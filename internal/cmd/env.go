package cmd

import (
	"github.com/urfave/cli/v3"
)

var envCmd = &cli.Command{
	Name:  "env",
	Usage: "Generate .env files for local development",
	Description: "Decrypts SOPS files and generates .env files in each app directory. " +
		"Resolves all environment variable types (value, resource, sops) for the development environment. " +
		"Use this for local development setup.",
}
