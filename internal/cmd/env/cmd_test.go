package env

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

func setup(t *testing.T, def *appdef.Definition) cmdtools.CommandInput {
	t.Helper()

	fs := afero.NewMemMapFs()
	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: def,
		Manifest:    manifest.NewTracker(),
	}

	return input
}
