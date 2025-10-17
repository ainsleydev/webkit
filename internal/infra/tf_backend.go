package infra

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/spf13/afero"
)

const (
	backendFileName = "backend.hcl"
)

// writeBackendConfig writes the Terraform backend configuration file
// using credentials from TFEnvironment.
func writeBackendConfig(fs afero.Fs, infraDir string, env TFEnvironment) error {
	content := fmt.Sprintf(`access_key = "%s"
secret_key = "%s"`, env.BackBlazeKeyID, env.BackBlazeApplicationKey)

	path := filepath.Join(infraDir, backendFileName)
	if err := afero.WriteFile(fs, path, []byte(content), 0600); err != nil {
		return errors.Wrap(err, "writing backend config")
	}

	return nil
}
