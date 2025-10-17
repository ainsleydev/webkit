package infra

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/pkg/env"
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

// writeS3Backend writes the complete Terraform backend configuration
// with a dynamic key based on project name and environment
func (t *Terraform) writeS3Backend(infraDir string, environment env.Environment) (string, error) {
	// Generate state file key: project-name/environment/terraform.tfstate
	stateKey := fmt.Sprintf("%s/%s/terraform.tfstate", t.appDef.Project.Name, environment)

	content := fmt.Sprintf(`bucket = "%s"
key                         = "%s"
region                      = "eu-central-003"
skip_credentials_validation = true
skip_region_validation      = true
skip_requesting_account_id  = true
use_path_style              = true
endpoint                    = "https://s3.eu-central-003.backblazeb2.com"
access_key                  = "%s"
secret_key                  = "%s"`,
		t.env.BackBlazeBucket,
		stateKey,
		t.env.BackBlazeKeyID,
		t.env.BackBlazeApplicationKey)

	path := filepath.Join(infraDir, backendFileName)
	if err := afero.WriteFile(t.fs, path, []byte(content), 0600); err != nil {
		return "", errors.Wrap(err, "writing backend config")
	}

	return path, nil
}
