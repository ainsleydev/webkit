package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
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

// getEnvironmentVars extracts vars for the specified environment.
func getEnvironmentVars(app appdef.Environment, environment env.Environment) (appdef.EnvVar, error) {
	switch environment {
	case env.Development:
		return app.Dev, nil
	case env.Staging:
		return app.Staging, nil
	case env.Production:
		return app.Production, nil
	default:
		return nil, fmt.Errorf("unsupported environment: %s", environment)
	}
}

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
