package secrets

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/executil"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestGet(t *testing.T) {
	t.Parallel()

	if !executil.Exists("sops") {
		t.Skip("sops CLI not found in PATH; skipping integration test")
	}

	tmpDir := t.TempDir()
	secretsDir := filepath.Join(tmpDir, secrets.FilePath)
	require.NoError(t, os.MkdirAll(secretsDir, os.ModePerm))

	envName := env.Production
	key := "SECRET_KEY"
	value := "supersecret"

	plainPath := filepath.Join(secretsDir, envName+".yaml")

	// Write plaintext YAML
	content := fmt.Sprintf("%s: %q\n", key, value)
	require.NoError(t, os.WriteFile(plainPath, []byte(content), os.ModePerm))

	// Encrypt using SOPS + age
	cmd := exec.CommandContext(t.Context(), "sops", "--encrypt", "--in-place", plainPath)
	out, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "failed to encrypt test file with sops: %s", string(out))

	// Prepare CLI command
	cmdCLI := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "env"},
			&cli.StringFlag{Name: "key"},
		},
	}
	require.NoError(t, cmdCLI.Set("env", envName))
	require.NoError(t, cmdCLI.Set("key", key))

	input := cmdtools.CommandInput{
		Command: cmdCLI,
		BaseDir: tmpDir,
	}

	// Run Get
	err = Get(t.Context(), input)
	require.NoError(t, err)
}
