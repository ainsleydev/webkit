package secrets

import (
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

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

func setupEncryptedProdFile(t *testing.T, content string) cmdtools.CommandInput {
	t.Helper()

	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)
	t.Setenv(age.KeyEnvVar, ageIdentity.String())

	tmpDir := t.TempDir()
	fs := afero.NewBasePathFs(afero.NewOsFs(), tmpDir)

	input := cmdtools.CommandInput{
		FS:      fs,
		BaseDir: tmpDir,
		Command: GetCmd,
	}
	input.Printer().SetWriter(io.Discard)
	err = CreateFiles(t.Context(), input)
	require.NoError(t, err)

	path := "resources/secrets/production.yaml"
	err = afero.WriteFile(input.FS, path, []byte(content), os.ModePerm)
	require.NoError(t, err)

	return input
}
