package env

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var SyncCmd = &cli.Command{
	Name:        "sync",
	Usage:       "Sync secrets to env files from app.json",
	Description: "Reads app.json and adds creates or updates .env files in the relevant app directories",
	Action:      cmdtools.Wrap(Sync),
}

func Sync(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	enviro := "production"

	//err := secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
	//	SOPSClient: input.SOPSClient(),
	//})
	//if err != nil {
	//	return err
	//}

	for _, app := range appDef.Apps {
		mergedApp, ok := appDef.MergeAppEnvironment(app.Name)
		if !ok {
			continue
		}

		envMap := make(map[string]string)
		for k, v := range mergedApp.Production {
			envMap[k] = cast.ToString(v.Value)
		}

		buf, err := godotenv.Marshal(envMap)
		if err != nil {
			return err
		}

		if err = input.FS.MkdirAll(app.Path, os.ModePerm); err != nil {
			return err
		}

		file := fmt.Sprintf(".env.%s", enviro)
		envPath := filepath.Join(app.Path, file)
		if err = afero.WriteFile(input.FS, envPath, []byte(buf), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
