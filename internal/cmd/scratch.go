package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
)

var scratchCmd = &cli.Command{
	Name:   "scratch",
	Hidden: true,
	Action: cmdtools.Wrap(func(ctx context.Context, input cmdtools.CommandInput) error {
		appDef := input.AppDef()
		cfg := secrets.ResolveConfig{SOPSClient: input.SOPSClient()}

		err := secrets.Resolve(ctx, appDef, cfg)
		if err != nil {
			return err
		}

		indent, err := json.MarshalIndent(cfg, "", "\t")
		if err != nil {
			return err
		}

		fmt.Println(string(indent))

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
