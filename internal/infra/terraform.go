package infra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/util/executil"
	tfembed "github.com/ainsleydev/webkit/platform/terraform"
)

type Terraform struct {
	path   string
	tmpDir string
	tf     *tfexec.Terraform
}

// NewTerraform creates a new Terraform manager by locating
// the terraform binary on the system.
//
// Returns an error if terraform cannot be found in PATH.
func NewTerraform(ctx context.Context) (*Terraform, error) {
	path, err := getTerraformPath(ctx)
	if err != nil {
		return nil, err
	}

	return &Terraform{
		path: path,
	}, nil
}

const tmpFolderPattern = "webkit-tf"

func (t *Terraform) Init(ctx context.Context) error {
	tmpDir, err := os.MkdirTemp("", tmpFolderPattern)
	if err != nil {
		return err
	}
	t.tmpDir = tmpDir

	err = fsext.CopyAllEmbed(tfembed.Templates, tmpDir)
	if err != nil {
		return err
	}

	tfDir := filepath.Join(tmpDir, "base")
	tf, err := tfexec.NewTerraform(tfDir, t.path)
	if err != nil {
		return fmt.Errorf("creating terraform executor: %w", err)
	}
	t.tf = tf

	if err = tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return fmt.Errorf("terraform init: %w", err)
	}
	return nil
}

func (t *Terraform) Plan(ctx context.Context) error {
	if t.tf == nil {
		return fmt.Errorf("terraform not initialized: call Init() first")
	}

	// TODO: Replace hardcoded vars with GenerateTFVars() output
	_, err := t.tf.Plan(ctx,
		tfexec.Var("project_name=hey!"),
		tfexec.Var("environment=production!"),
	)
	if err != nil {
		return fmt.Errorf("executing plan: %w", err)
	}

	return nil
}

// Show retrieves the current Terraform state as a top level
// representation of the Terraform state.
func (t *Terraform) Show(ctx context.Context) (*tfjson.State, error) {
	if t.tf == nil {
		return nil, fmt.Errorf("terraform not initialized: call Init() first")
	}
	return t.tf.Show(ctx)
}

func (t *Terraform) Cleanup() {
	if t.tmpDir != "" {
		_ = os.RemoveAll(t.tmpDir) //nolint
	}
}

func getTerraformPath(ctx context.Context) (string, error) {
	whichCmd := executil.NewCommand("which", "terraform")
	run, err := executil.DefaultRunner().Run(ctx, whichCmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(run.Output), nil
}
