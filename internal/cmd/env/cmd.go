package env

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

	buf, err := dotEnvMarshaller(envMap)
	if err != nil {
		return err
	}

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
	err = args.Input.FS.MkdirAll(outputDir, os.ModePerm)
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

	provider := make(secrets.TerraformOutputProvider)
	for resourceName, outputs := range result.Resources {
		for outputName, value := range outputs {
			key := secrets.OutputKey{
				Environment:  environment,
				ResourceName: resourceName,
				OutputName:   outputName,
			}
			provider[key] = value
		}
	}

	return &provider, nil
}
