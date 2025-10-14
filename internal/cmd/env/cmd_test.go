package env

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

func setup(t *testing.T, def *appdef.Definition) (cmdtools.CommandInput, *bytes.Buffer) {
	t.Helper()

	fs := afero.NewMemMapFs()
	buf := &bytes.Buffer{}
	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: def,
	}
	input.Printer().SetWriter(buf)

	return input, buf
}
