package secrets

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/pkg/env"
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

func setupEncryptedProdFile(t *testing.T, content string) (cmdtools.CommandInput, *bytes.Buffer) {
	t.Helper()

	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)
	t.Setenv(age.KeyEnvVar, ageIdentity.String())

	tmpDir := t.TempDir()
	fs := afero.NewBasePathFs(afero.NewOsFs(), tmpDir)

	buf := &bytes.Buffer{}
	input := cmdtools.CommandInput{
		FS:      fs,
		BaseDir: tmpDir,
		Command: GetCmd,
	}
	input.Printer().SetWriter(buf)
	err = Scaffold(t.Context(), input)
	require.NoError(t, err)

	path := secrets.FilePathFromEnv(env.Production)
	err = afero.WriteFile(input.FS, path, []byte(content), os.ModePerm)
	require.NoError(t, err)

	return input, buf
}
