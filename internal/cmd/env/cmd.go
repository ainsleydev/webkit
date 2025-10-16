package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
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
	},
}

// environmentsWithDotEnv defines the environments to generate .env
// files for.
var environmentsWithDotEnv = []env.Environment{
	env.Development,
	env.Production,
}

type writeArgs struct {
	FS          afero.Fs
	Vars        appdef.EnvVar
	App         appdef.App
	Environment env.Environment
	IsScaffold  bool
}

var dotEnvMarshaller = godotenv.Marshal

// writeMapToFile writes environment variables to dotenv file.
func writeMapToFile(args writeArgs) error {
	gen := scaffold.New(args.FS)

	envMap := make(map[string]string)
	for k, v := range args.Vars {
		envMap[k] = cast.ToString(v.Value)
	}

	buf, err := dotEnvMarshaller(envMap)
	if err != nil {
		return err
	}

	err = args.FS.MkdirAll(args.App.Path, os.ModePerm)
	if err != nil {
		return err
	}

	file := ".env"
	if args.Environment != env.Development {
		file = fmt.Sprintf(".env.%s", args.Environment)
	}

	envPath := filepath.Join(args.App.Path, file)

	var opts []scaffold.Option
	opts = append(opts, scaffold.WithNotice(true))
	if args.IsScaffold {
		opts = append(opts, scaffold.WithScaffoldMode())
	}

	return gen.Bytes(envPath, []byte(buf), opts...)
}
