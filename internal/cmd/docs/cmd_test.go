package docs

import (
	"io"
	"testing"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

func setup(t *testing.T, fs afero.Fs, appDef *appdef.Definition) cmdtools.CommandInput { //nolint
	t.Helper()

	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: appDef,
		Manifest:    manifest.NewTracker(),
	}
	input.Printer().SetWriter(io.Discard)

	return input
}
