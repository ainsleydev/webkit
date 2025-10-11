package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/lipgloss/tree"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations/infra"
	"github.com/ainsleydev/webkit/internal/printer"
)

var scratchCmd = &cli.Command{
	Name:   "scratch",
	Hidden: true,
	Action: cmdtools.Wrap(func(ctx context.Context, input cmdtools.CommandInput) error {

		c := printer.New(os.Stdout)

		c.Success("Deployment completed successfully.")
		c.Info("Fetching configuration...")
		c.Warn("Config file deprecated, using fallback.")
		c.Error("Failed to connect to database.")

		c.Table(
			[]string{"Service", "Status"},
			[][]string{
				{"Database", "Running"},
				{"API", "Stopped"},
			},
		)

		c.List("Config loaded", "Deployment started", "Done ✅")

		c.Tree("project",
			"cmd",
			"internal",
			tree.New().Root("pkg").Child("config", "printer", "infra"),
		)

		//p.Title("Hello")
		//p.LineBreak()
		//
		//p.FatalError(errors.New("Eree"))
		//p.Title("World")
		//p.StatusList("hey", []printer.StatusListItem{
		//	{false, "Status"},
		//})
		//prov, err := age.NewProvider()
		//if err != nil {
		//	return err
		//}
		//
		//client := sops.NewClient(prov)
		//
		//path := filepath.Join(input.BaseDir, secrets.FilePath, "production.yaml")
		//toMap, err := sops.DecryptFileToMap(client, path)
		//
		//fmt.Print(toMap, err)

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
