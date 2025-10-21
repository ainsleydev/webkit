package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

func setup(t *testing.T, fs afero.Fs, appDef *appdef.Definition) cmdtools.CommandInput {
	t.Helper()

	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: appDef,
		Manifest:    manifest.NewTracker(),
	}
	input.Printer().SetWriter(io.Discard)

	return input
}

func setupWithPrinter(t *testing.T, fs afero.Fs, def *appdef.Definition) (cmdtools.CommandInput, *bytes.Buffer) {
	t.Helper()

	input := setup(t, fs, def)
	buf := &bytes.Buffer{}
	input.Printer().SetWriter(buf)

	return input, buf
}
