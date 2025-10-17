package infra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-json"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/util/executil"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/platform/terraform"
)

// Terraform represents the type for interacting with the
// terraform exec CLI.
type Terraform struct {
	path   string
	tmpDir string
	env    TFEnvironment
	tf     terraformExecutor
	fs     afero.Fs
}

//go:generate go tool go.uber.org/mock/mockgen -source=tf.go -destination ./internal/tfmocks/tf.go -package=tfmocks

// terraformExecutor defines the interface for terraform operations
// using tf exec.
type terraformExecutor interface {
	Init(ctx context.Context, opts ...tfexec.InitOption) error
	Plan(ctx context.Context, opts ...tfexec.PlanOption) (bool, error)
	Apply(ctx context.Context, opts ...tfexec.ApplyOption) error
	ShowPlanFileRaw(ctx context.Context, planPath string, opts ...tfexec.ShowOption) (string, error)
	ShowPlanFile(ctx context.Context, planPath string, opts ...tfexec.ShowOption) (*tfjson.Plan, error)
}

// NewTerraform creates a new Terraform manager by locating
// the terraform binary on the system.
//
// Returns an error if terraform cannot be found in PATH.
func NewTerraform(ctx context.Context, fs afero.Fs) (*Terraform, error) {
	path, err := getTerraformPath(ctx)
	if err != nil {
		return nil, err
	}
	tfEnv, err := ParseTFEnvironment()
	if err != nil {
		return nil, err
	}
	return &Terraform{
		path: path,
		fs:   fs,
		env:  tfEnv,
	}, nil
}

const tmpFolderPattern = "webkit-tf"

// Init initialises the WebKit terraform provider by copying all the
// terraform embedded templates to a temporary directory on the
// filesystem.
//
// Backend configuration and provider config are also written as
// part of this process.
//
// Must be called before Plan() or Apply()
func (t *Terraform) Init(ctx context.Context) error {
	tmpDir, err := afero.TempDir(t.fs, "", tmpFolderPattern)
	if err != nil {
		return errors.Wrap(err, "creating tf tmp dir")
	}
	t.tmpDir = tmpDir

	err = fsext.CopyAllEmbed(tfembed.Templates, tmpDir)
	if err != nil {
		return err
	}

	tfDir := filepath.Join(tmpDir, "base")
	tf, err := tfexec.NewTerraform(tfDir, t.path)
	if err != nil {
		return errors.Wrap(err, "creating terraform executor")
	}
	t.tf = tf

	baseDir := filepath.Join(tmpDir, "base")
	if err = writeBackendConfig(t.fs, baseDir, t.env); err != nil {
		return err
	}

	if err = writeProviderConfig(t.fs, baseDir, t.env); err != nil {
		return err
	}

	if err = tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return errors.Wrap(err, "initialising tf")
	}

	return nil
}

// PlanOutput is the result of calling Plan.
type PlanOutput struct {
	// Human-readable output (the output that's usually
	// in the terminal when running terraform plan).
	Output string

	// The JSON contents of the plan for more of a
	// detailed look.
	Plan *tfjson.Plan
}

// Plan generates a Terraform execution plan showing what actions Terraform
// will take to reach the desired state defined in the definition.
//
// Must be called after Init().
func (t *Terraform) Plan(ctx context.Context, env env.Environment, def *appdef.Definition) (PlanOutput, error) {
	if err := t.prepareVars(env, def); err != nil {
		return PlanOutput{}, err
	}

	planFilePath := filepath.Join(t.tmpDir, "base", "plan.tfplan")
	_, err := t.tf.Plan(ctx, tfexec.Out(planFilePath))
	if err != nil {
		return PlanOutput{}, fmt.Errorf("terraform plan failed: %w", err)
	}

	// Human-readable output.
	output, err := t.tf.ShowPlanFileRaw(ctx, planFilePath)
	if err != nil {
		return PlanOutput{}, fmt.Errorf("showing plan file: %w", err)
	}

	// Fully typed output.
	file, err := t.tf.ShowPlanFile(ctx, planFilePath)
	if err != nil {
		return PlanOutput{}, fmt.Errorf("showing plan file: %w", err)
	}

	return PlanOutput{
		Output: output,
		Plan:   file,
	}, nil
}

// Apply executes terraform apply to provision infrastructure based on
// the app definition provided.
//
// Must be called after Init().
func (t *Terraform) Apply(ctx context.Context, env env.Environment, def *appdef.Definition) error {
	if err := t.prepareVars(env, def); err != nil {
		return err
	}

	if err := t.tf.Apply(ctx); err != nil {
		return fmt.Errorf("terraform apply failed: %w", err)
	}

	return nil
}

// Cleanup removes all the temporary directories that we're
// created during the terraform init process.
//
// Ideally should be called after Init().
func (t *Terraform) Cleanup() {
	if t.tmpDir != "" {
		_ = os.RemoveAll(t.tmpDir) //nolint
	}
}

func (t *Terraform) hasInitialised() error {
	if t.tf == nil {
		return errors.New("terraform not initialized: call Init() first")
	}
	return nil
}

func (t *Terraform) prepareVars(env env.Environment, def *appdef.Definition) error {
	if err := t.hasInitialised(); err != nil {
		return err
	}

	vars, err := tfVarsFromDefinition(env, def)
	if err != nil {
		return errors.Wrap(err, "generating terraform variables")
	}

	if err = t.writeTFVarsFile(vars); err != nil {
		return errors.Wrap(err, "writing tfvars file")
	}

	return nil
}

func getTerraformPath(ctx context.Context) (string, error) {
	whichCmd := executil.NewCommand("which", "terraform")
	run, err := executil.DefaultRunner().Run(ctx, whichCmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(run.Output), nil
}
