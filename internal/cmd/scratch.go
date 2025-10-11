package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations/infra"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
)

var scratchCmd = &cli.Command{
	Name:   "scratch",
	Hidden: true,
	Action: cmdtools.Wrap(func(ctx context.Context, input cmdtools.CommandInput) error {

		prov, err := age.NewProvider()
		if err != nil {
			return err
		}

		client := sops.NewClient(prov)

		path := filepath.Join(input.BaseDir, secrets.FilePath, "production.yaml")
		toMap, err := sops.DecryptFileToMap(client, path)

		fmt.Print(toMap, err)

		return nil

		return infra.Test(ctx, input)
		return nil

		//// Create reflector with custom configuration
		//reflector := &jsonschema.Reflector{
		//	AllowAdditionalProperties: false,
		//	DoNotReference:            false,
		//	ExpandedStruct:            true,
		//}
		//
		//// Generate schema from Definition struct
		//schema := reflector.Reflect(&appdef.Definition{})
		//
		//// Add metadata
		//schema.Title = "WebKit Application Manifest"
		//schema.Description = "Schema for webkit app.json configuration file"
		//schema.Version = "1.0.0"
		//
		//data, err := json.MarshalIndent(schema, "", "  ")
		//if err != nil {
		//	return err
		//}
		//
		//// Write to file
		//return os.WriteFile("schema-test-2.json", data, 0644)

		//input.AppDef()
		//
		//reflector := jsonschema.Reflector{}
		//
		//schema, err := reflector.Reflect(appdef.Definition{})
		//if err != nil {
		//	return err
		//}
		//
		//j, err := json.MarshalIndent(schema, "", " ")
		//if err != nil {
		//	return err
		//}
		//
		//return os.WriteFile("schema-test.json", j, 0600)
	}),
}
