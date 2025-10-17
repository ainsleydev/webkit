package infra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-json"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/util/executil"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/platform/terraform"
)

type Terraform struct {
	path   string
	tmpDir string
	tf     *tfexec.Terraform
	fs     afero.Fs
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
		fs:   afero.NewOsFs(),
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

	tfEnv, err := ParseTFEnvironment()
	if err != nil {
		return fmt.Errorf("parsing terraform environment: %w", err)
	}

	baseDir := filepath.Join(tmpDir, "base")
	if err = writeBackendConfig(t.fs, baseDir, tfEnv); err != nil {
		return fmt.Errorf("writing backend config: %w", err)
	}

	if err = writeProviderConfig(t.fs, baseDir, tfEnv); err != nil {
		return fmt.Errorf("writing provider config: %w", err)
	}

	if err = tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return fmt.Errorf("terraform init: %w", err)
	}

	return nil
}

func (t *Terraform) Plan(ctx context.Context, env env.Environment, def *appdef.Definition) (*tfjson.Plan, error) {
	if t.tf == nil {
		return nil, fmt.Errorf("terraform not initialized: call Init() first")
	}

	vars, err := tfVarsFromDefinition(env, def)
	if err != nil {
		return nil, fmt.Errorf("generating terraform variables: %w", err)
	}

	if err = t.writeTFVarsFile(vars); err != nil {
		return nil, fmt.Errorf("writing tfvars file: %w", err)
	}

	planFilePath := filepath.Join(t.tmpDir, "base", "plan.tfplan")
	_, err = t.tf.Plan(ctx, tfexec.Out(planFilePath))
	if err != nil {
		return nil, fmt.Errorf("terraform plan failed: %w", err)
	}

	// Get the human-readable plan output
	planOutput, err := t.tf.ShowPlanFileRaw(ctx, planFilePath)
	if err != nil {
		return nil, fmt.Errorf("showing plan file: %w", err)
	}

	fmt.Println(string(planOutput))

	return t.tf.ShowPlanFile(ctx, planFilePath)
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
