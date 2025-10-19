package cmd

import (
	"bytes"
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

func setupWithPrinter(t *testing.T, fs afero.Fs, def *appdef.Definition) (cmdtools.CommandInput, *bytes.Buffer) {
	t.Helper()

	input := setup(t, fs, def)
	buf := &bytes.Buffer{}
	input.Printer().SetWriter(buf)

	return input, buf
}
