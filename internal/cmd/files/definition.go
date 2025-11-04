package files

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/appdef/internal/jsonformat"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/version"
)

// Definition updates the app.json file with the current CLI version.
// This ensures the webkit_version field stays in sync with the installed CLI.
func Definition(_ context.Context, input cmdtools.CommandInput) error {
	def := input.AppDef()

	// Update the webkit_version to match the current CLI version.
	def.WebkitVersion = version.Version

	data, err := identMarshaller(def, "", "\t")
	if err != nil {
		return errors.Wrap(err, "marshaling definition")
	}

	// Apply custom formatting to inline specific objects.
	data, err = jsonformat.Format(data)
	if err != nil {
		return errors.Wrap(err, "formatting JSON")
	}

	// Add trailing newline for better git diffs.
	data = append(data, '\n')

	if err = afero.WriteFile(input.FS, appdef.JsonFileName, data, 0o644); err != nil {
		return errors.Wrap(err, "writing app.json")
	}

	return nil
}
