// Package env provides commands for managing environment variable files (.env)
// based on the application definition.
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

// Command defines the environment commands for generating and managing
// environment variable files.
var Command = &cli.Command{
	Name:        "env",
	Usage:       "Manage environment variables",
	Description: "Command for working with the environment files defined in app.json",
	Commands: []*cli.Command{
		ScaffoldCmd,
		SyncCmd,
	},
}

// environmentsWithDotEnv defines the environments to generate .env
// files for.
var environmentsWithDotEnv = []env.Environment{
	env.Development,
	env.Production,
}

// writeArgs contains the parameters needed to write environment variables to a file.
type writeArgs struct {
	Input       cmdtools.CommandInput
	Vars        appdef.EnvVar
	App         appdef.App
	Environment env.Environment
	IsScaffold  bool
}

// dotEnvMarshaller is the function used to marshal environment variables to dotenv format.
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

	err = args.Input.FS.MkdirAll(args.App.Path, os.ModePerm)
	if err != nil {
		return err
	}

	file := ".env"
	if args.Environment != env.Development {
		file = fmt.Sprintf(".env.%s", args.Environment)
	}

	envPath := filepath.Join(args.App.Path, file)

	var opts []scaffold.Option
	if args.IsScaffold {
		opts = append(opts, scaffold.WithScaffoldMode())
	}
	opts = append(opts, scaffold.WithTracking(manifest.SourceProject()))

	return args.Input.Generator().Bytes(envPath, []byte(buf), opts...)
}
