package infra

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	mockinfra "github.com/ainsleydev/webkit/internal/infra/mocks"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/secrets/age"
)

var mtx sync.Mutex

func setup(t *testing.T, appDef *appdef.Definition, mock *mockinfra.MockManager, initErr bool) (cmdtools.CommandInput, func()) {
	t.Helper()

	mtx.Lock()
	defer mtx.Unlock()

	var err error
	if initErr {
		err = errors.New("init error")
	}

	mock.EXPECT().
		Init(t.Context()).
		Return(err)

	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)
	t.Setenv(age.KeyEnvVar, ageIdentity.String())

	fs := afero.NewMemMapFs()
	input := cmdtools.CommandInput{
		FS:          fs,
		BaseDir:     t.TempDir(),
		AppDefCache: appDef,
		Manifest:    manifest.NewTracker(),
		Command:     &cli.Command{},
	}
	input.Printer().SetWriter(io.Discard)

	orig := newTerraform
	newTerraform = func(_ context.Context, _ *appdef.Definition, _ *manifest.Tracker) (infra.Manager, error) {
		return mock, nil
	}

	return input, func() {
		newTerraform = orig
	}
}

func setupWithPrinter(t *testing.T, def *appdef.Definition, manager *mockinfra.MockManager, initError bool) (cmdtools.CommandInput, *bytes.Buffer, func()) { //nolint
	t.Helper()

	input, teardown := setup(t, def, manager, initError)
	buf := &bytes.Buffer{}
	input.Printer().SetWriter(buf)

	return input, buf, teardown
}
