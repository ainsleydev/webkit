package payload

import (
	"io"
	"testing"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/util/executil"
)

func setup(t *testing.T) (afero.Fs, cmdtools.CommandInput) {
	t.Helper()

	fs := afero.NewMemMapFs()
	input := cmdtools.CommandInput{
		FS:      fs,
		Command: &cli.Command{},
	}
	input.Printer().SetWriter(io.Discard)

	return fs, input
}

func setupWithRunner(t *testing.T, runner executil.Runner) (afero.Fs, cmdtools.CommandInput) {
	t.Helper()

	fs, input := setup(t)
	input.Runner = runner

	return fs, input
}
