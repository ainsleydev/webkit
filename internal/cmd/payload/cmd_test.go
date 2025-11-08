package payload

import (
	"io"
	"testing"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/cmdtools"
)

func setup(t *testing.T) (afero.Fs, cmdtools.CommandInput) {
	t.Helper()

	fs := afero.NewMemMapFs()
	input := cmdtools.CommandInput{
		FS: fs,
	}
	input.Printer().SetWriter(io.Discard)

	return fs, input
}
