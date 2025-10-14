package env

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/mocks"
)

func setup(t *testing.T, def *appdef.Definition) (cmdtools.CommandInput, *bytes.Buffer) {
	t.Helper()

	h := mocks.NewMockTester(gomock.NewController(t))

	fs := afero.NewMemMapFs()
	buf := &bytes.Buffer{}
	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: def,
	}
	input.Printer().SetWriter(buf)

	return input, buf
}
