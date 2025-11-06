package env

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

// Command defines the env commands for interacting and generating
// env file artifacts.
var Command = &cli.Command{
	Name:        "env",
	Usage:       "Manage environment variables",
	Description: "Command for working with the environment files defined in app.json",
	Commands: []*cli.Command{
		ScaffoldCmd,
		SyncCmd,
		GenerateCmd,
	},
}

// environmentsWithDotEnv defines the environments to generate .env
// files for.
var environmentsWithDotEnv = []env.Environment{
	env.Development,
	env.Production,
}

type writeArgs struct {
	Input            cmdtools.CommandInput
	Vars             appdef.EnvVar
	App              appdef.App
	Environment      env.Environment
	IsScaffold       bool
	CustomOutputPath string
}

var dotEnvMarshaller = godotenv.Marshal

// writeMapToFile writes environment variables to dotenv file.
func writeMapToFile(args writeArgs) error {
	envMap := make(map[string]string)
	for k, v := range args.Vars {
		envMap[k] = cast.ToString(v.Value)
	}

	// Use custom marshaller that doesn't quote unnecessarily
	// This prevents Docker Swarm from including quotes in the actual env var values
	buf := marshalEnvWithoutQuotes(envMap)

	var envPath string
	if args.CustomOutputPath != "" {
		envPath = args.CustomOutputPath
	} else {
		file := ".env"
		if args.Environment != env.Development {
			file = fmt.Sprintf(".env.%s", args.Environment)
		}
		envPath = filepath.Join(args.App.Path, file)
	}

	outputDir := filepath.Dir(envPath)
	err := args.Input.FS.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	var opts []scaffold.Option
	if args.IsScaffold {
		opts = append(opts, scaffold.WithScaffoldMode())
	}
	opts = append(opts, scaffold.WithTracking(manifest.SourceProject()))

	return args.Input.Generator().Bytes(envPath, []byte(buf), opts...)
}

// envSuffix returns the .env file suffix for an environment.
func envSuffix(environment env.Environment) string {
	if environment == env.Development {
		return ""
	}
	return fmt.Sprintf(".%s", environment)
}

// fetchTerraformOutputs fetches Terraform outputs for the specified environment.
// Returns a TerraformOutputProvider containing resource outputs.
//
// This function manages the full Terraform lifecycle (create, init, output, cleanup).
// See also: infra/cmd.go:fetchTerraformOutputs for a similar function that uses
// an existing Terraform instance.
func fetchTerraformOutputs(
	ctx context.Context,
	input cmdtools.CommandInput,
	environment env.Environment,
) (*secrets.TerraformOutputProvider, error) {
	tf, err := infra.NewTerraform(ctx, input.AppDef(), input.Manifest)
	if err != nil {
		return nil, errors.Wrap(err, "creating terraform manager")
	}

	if err := tf.Init(ctx); err != nil {
		return nil, errors.Wrap(err, "initialising terraform")
	}
	defer tf.Cleanup()

	result, err := tf.Output(ctx, environment)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving terraform outputs")
	}

	provider := secrets.TransformOutputs(result, environment)
	return &provider, nil
}

// marshalEnvWithoutQuotes marshals environment variables without adding quotes.
// This is necessary for Docker Swarm env_files which don't strip quotes like docker-compose does.
// Only adds quotes when the value contains spaces, newlines, or is empty.
func marshalEnvWithoutQuotes(envMap map[string]string) string {
	var builder strings.Builder
	for key, value := range envMap {
		// Only quote if value contains spaces, newlines, or is empty
		// Docker Swarm env_files doesn't handle quotes like docker-compose
		if strings.ContainsAny(value, " \n\t") || value == "" {
			// Escape any quotes in the value
			escapedValue := strings.ReplaceAll(value, `"`, `\"`)
			builder.WriteString(fmt.Sprintf("%s=\"%s\"\n", key, escapedValue))
		} else {
			builder.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		}
	}
	return builder.String()
}
