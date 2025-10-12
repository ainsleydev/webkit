package secrets

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
)

type mockEncrypterDecrypter struct {
	err error
}

func (m mockEncrypterDecrypter) Encrypt(_ string) error {
	return m.err
}

func (m mockEncrypterDecrypter) Decrypt(_ string) error {
	return m.err
}

var _ sops.EncrypterDecrypter = (*mockEncrypterDecrypter)(nil)

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
	err = CreateFiles(t.Context(), input)
	require.NoError(t, err)

	path := "resources/secrets/production.yaml"
	err = afero.WriteFile(input.FS, path, []byte(content), os.ModePerm)
	require.NoError(t, err)

	return input, buf
}
