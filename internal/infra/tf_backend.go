package infra

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
	"github.com/ainsleydev/webkit/pkg/env"
)

const (
	backendTfFileName  = "backend.tf"
	backendHclFileName = "backend.hcl"
)

// writeS3Backend writes the complete Terraform backend configuration
// with a dynamic key based on project name and environment
func (t *Terraform) writeS3Backend(infraDir string, environment env.Environment) (string, error) {
	gen := scaffold.New(t.fs, nil)
	gen.Printer.SetWriter(io.Discard)

	// For example, project-name/environment/terraform.tfstate
	folder := fmt.Sprintf("%s/%s/terraform.tfstate", t.appDef.Project.Name, environment)

	data := map[string]any{
		"Bucket":    t.env.BackBlazeBucket,
		"Key":       folder,
		"AccessKey": t.env.BackBlazeKeyID,
		"SecretKey": t.env.BackBlazeApplicationKey,
	}

	// Generate the .hcl file which contains all the sensitive data.
	backendHclPath := filepath.Join(infraDir, backendHclFileName)
	err := gen.Template(
		backendHclPath,
		templates.MustLoadTemplate("terraform/backend.hcl.tmpl"),
		data,
	)
	if err != nil {
		return "", err
	}

	// Generate the backend.tf file which just indicates to Terraform
	// that it should use S3.
	return backendHclPath, gen.Template(
		filepath.Join(infraDir, backendTfFileName),
		templates.MustLoadTemplate("terraform/backend.tf"),
		nil,
	)
}
