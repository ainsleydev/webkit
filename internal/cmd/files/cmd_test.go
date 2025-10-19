package files

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

func setup(t *testing.T, fs afero.Fs, appDef *appdef.Definition) cmdtools.CommandInput {
	t.Helper()

	return cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: appDef,
		Manifest:    manifest.NewTracker(),
	}
}
