package infra

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/spf13/afero"
)

const (
	backendFileName  = "backend.hcl"
	providerFileName = "providers.auto.tfvars"
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

// writeProviderConfig writes the Terraform provider configuration file
// with DigitalOcean credentials from TFEnvironment.
func writeProviderConfig(fs afero.Fs, infraDir string, env TFEnvironment) error {
	content := fmt.Sprintf(`do_token = "%s"
spaces_access_id  = "%s"
spaces_secret_key = "%s"`,
		env.DigitalOceanAPIKey, env.DigitalOceanSpacesAccessKey, env.DigitalOceanSpacesSecretKey)

	path := filepath.Join(infraDir, providerFileName)
	if err := afero.WriteFile(fs, path, []byte(content), 0600); err != nil {
		return errors.Wrap(err, "writing provider config")
	}

	return nil
}
