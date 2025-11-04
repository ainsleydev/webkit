package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-json"
	"github.com/kaptinlin/messageformat-go/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/ghapi"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/util/executil"
	"github.com/ainsleydev/webkit/pkg/enforce"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/platform/terraform"
)

// Terraform represents the type for interacting with the
// terraform exec CLI.
type Terraform struct {
	appDef          *appdef.Definition
	path            string
	tmpDir          string
	env             TFEnvironment
	tf              terraformExecutor
	manifest        *manifest.Tracker
	fs              afero.Fs
	ghClient        ghapi.Client
	useLocalBackend bool
}

//go:generate go tool go.uber.org/mock/mockgen -source=tf.go -destination ./internal/tfmocks/tf.go -package=tfmocks

// terraformExecutor defines the interface for terraform operations
// using tf exec.
type terraformExecutor interface {
	SetStdout(w io.Writer)
	SetStderr(w io.Writer)
	Init(ctx context.Context, opts ...tfexec.InitOption) error
	Plan(ctx context.Context, opts ...tfexec.PlanOption) (bool, error)
	Apply(ctx context.Context, opts ...tfexec.ApplyOption) error
	Destroy(ctx context.Context, opts ...tfexec.DestroyOption) error
	Output(ctx context.Context, opts ...tfexec.OutputOption) (map[string]tfexec.OutputMeta, error)
	Import(ctx context.Context, address string, id string, opts ...tfexec.ImportOption) error
	ShowPlanFileRaw(ctx context.Context, planPath string, opts ...tfexec.ShowOption) (string, error)
	ShowPlanFile(ctx context.Context, planPath string, opts ...tfexec.ShowOption) (*tfjson.Plan, error)
}

// NewTerraform creates a new Terraform manager by locating
// the terraform binary on the system.
//
// Returns an error if terraform cannot be found in PATH.
func NewTerraform(ctx context.Context, appDef *appdef.Definition, manifest *manifest.Tracker) (*Terraform, error) {
	enforce.NotNil(appDef, "app definition is required")
	enforce.NotNil(manifest, "manifest definition is required")

	path, err := getTerraformPath(ctx)
	if err != nil {
		return nil, err
	}

	tfEnv, err := ParseTFEnvironment()
	if err != nil {
		return nil, err
	}

	return &Terraform{
		appDef:          appDef,
		path:            path,
		fs:              afero.NewOsFs(),
		env:             tfEnv,
		ghClient:        ghapi.New(tfEnv.GithubTokenClassic),
		useLocalBackend: false,
		manifest:        manifest,
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

	initOpts := []tfexec.InitOption{
		tfexec.Upgrade(true),
		tfexec.Backend(!t.useLocalBackend), // Only use backend on prod.
		tfexec.Reconfigure(true),
	}

	if !t.useLocalBackend {
		backendPath, err := t.writeS3Backend(tfDir, env.Production)
		if err != nil {
			return err
		}
		initOpts = append(initOpts, tfexec.BackendConfig(backendPath))
	}

	if err = tf.Init(ctx, initOpts...); err != nil {
		return errors.Wrap(err, "initialising tf")
	}

	return nil
}

// PlanOutput is the result of calling Plan.
type PlanOutput struct {
	// Determines if there has been any changes to
	// the plan since running last.
	HasChanges bool

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
func (t *Terraform) Plan(ctx context.Context, env env.Environment) (PlanOutput, error) {
	if err := t.prepareVars(ctx, env); err != nil {
		return PlanOutput{}, err
	}

	planFilePath := filepath.Join(t.tmpDir, "base", "plan.tfplan")

	var vars []tfexec.PlanOption
	vars = append(vars, tfexec.Out(planFilePath))
	for _, v := range t.env.varStrings() {
		vars = append(vars, tfexec.Var(v))
	}

	changes, err := t.tf.Plan(ctx, vars...)
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
		HasChanges: changes,
		Output:     output,
		Plan:       file,
	}, nil
}

// ApplyOutput is the result of calling Apply.
type ApplyOutput struct {
	// Human-readable output (the output that's usually
	// in the terminal when running terraform apply).
	Output string
}

// Apply executes terraform apply to provision infrastructure based on
// the app definition provided.
//
// Must be called after Init().
func (t *Terraform) Apply(ctx context.Context, env env.Environment) (ApplyOutput, error) {
	if err := t.prepareVars(ctx, env); err != nil {
		return ApplyOutput{}, err
	}

	var outputBuf strings.Builder
	t.tf.SetStdout(&outputBuf)
	t.tf.SetStderr(&outputBuf)

	var vars []tfexec.ApplyOption
	for _, v := range t.env.varStrings() {
		vars = append(vars, tfexec.Var(v))
	}

	if err := t.tf.Apply(ctx, vars...); err != nil {
		return ApplyOutput{
			Output: outputBuf.String(),
		}, fmt.Errorf("terraform apply failed: %w", err)
	}

	return ApplyOutput{
		Output: outputBuf.String(),
	}, nil
}

// DestroyOutput is the result of calling Destroy.
type DestroyOutput struct {
	// Human-readable output (the output that's usually
	// in the terminal when running terraform destroy).
	Output string
}

// Destroy executes terraform destroy to tear down infrastructure
// based on the app definition provided.
//
// Must be called after Init().
func (t *Terraform) Destroy(ctx context.Context, env env.Environment) (DestroyOutput, error) {
	if err := t.prepareVars(ctx, env); err != nil {
		return DestroyOutput{}, err
	}

	var outputBuf strings.Builder
	t.tf.SetStdout(&outputBuf)
	t.tf.SetStderr(&outputBuf)

	var vars []tfexec.DestroyOption
	for _, v := range t.env.varStrings() {
		vars = append(vars, tfexec.Var(v))
	}

	if err := t.tf.Destroy(ctx, vars...); err != nil {
		return DestroyOutput{
			Output: outputBuf.String(),
		}, fmt.Errorf("terraform destroy failed: %w", err)
	}

	return DestroyOutput{
		Output: outputBuf.String(),
	}, nil
}

// OutputResult is the result of calling Output.
// It contains a structured map of all Terraform outputs.
// See platform/terraform/base/outputs.tf for spercifics.
//
//   - Resources: Maps resource names to their fields and values.
//     Example:
//     "store": {
//     "bucket_name": "my-website-store-temp",
//     "bucket_url": "my-website-store-temp.nyc3.digitaloceanspaces.com",
//     }
//
//   - Apps: Maps app names to their fields and values.
//     Example (empty if no apps provisioned):
//     "web-app": {
//     "app_url": "https://web-app.example.com",
//     "platform_provider": "digitalocean"
//     }
//
//   - Extra: Contains all other outputs that don’t fit into Resources or Apps.
type OutputResult struct {
	Resources map[string]map[string]any `json:"resources"`
	Apps      map[string]map[string]any `json:"apps"`
	Extra     map[string]any            `json:"extra"`
}

// Output retrieves all Terraform outputs for the specified environment.
// This reads the current terraform state and returns all output values.
//
// Must be called after Init().
func (t *Terraform) Output(ctx context.Context, env env.Environment) (OutputResult, error) {
	if err := t.prepareVars(ctx, env); err != nil {
		return OutputResult{}, err
	}

	rawOutputs, err := t.tf.Output(ctx)
	if err != nil {
		return OutputResult{}, errors.Wrap(err, "terraform output failed")
	}

	result := OutputResult{
		Resources: make(map[string]map[string]any),
		Apps:      make(map[string]map[string]any),
		Extra:     make(map[string]any),
	}

	// --- Resources ---
	if r, ok := rawOutputs["resources"]; ok {
		if err = json.Unmarshal(r.Value, &result.Resources); err != nil {
			return OutputResult{}, errors.Wrap(err, "unmarshalling resources")
		}
	}

	// --- Apps ---
	if r, ok := rawOutputs["apps"]; ok {
		if err = json.Unmarshal(r.Value, &result.Apps); err != nil {
			return OutputResult{}, errors.Wrap(err, "unmarshalling apps")
		}
	}

	// --- Extra (everything else) ---
	for key, meta := range rawOutputs {
		if key == "resources" || key == "apps" {
			continue
		}
		var val any
		if err = json.Unmarshal(meta.Value, &val); err != nil {
			val = string(meta.Value)
		}
		result.Extra[key] = val
	}

	return result, nil
}

type (
	// ImportInput contains the configuration for importing existing resources or apps.
	ImportInput struct {
		// ResourceName is the name of the resource or app in app.json.
		ResourceName string

		// ResourceID is the provider-specific ID (e.g., DigitalOcean cluster ID or app ID).
		ResourceID string

		// Environment specifies which environment to import into.
		Environment env.Environment

		// IsApp indicates whether we're importing an app (true) or a resource (false).
		IsApp bool
	}
	// ImportOutput contains the results of an import operation.
	ImportOutput struct {
		// ImportedResources lists the Terraform addresses that were imported.
		ImportedResources []string

		// Output contains the human-readable output from the import operations.
		Output string
	}
)

// Import imports an existing infrastructure resource or app into the Terraform state.
// This allows webkit to manage resources/apps that were created manually or outside of Terraform.
//
// Must be called after Init().
func (t *Terraform) Import(ctx context.Context, input ImportInput) (ImportOutput, error) {
	if err := t.prepareVars(ctx, input.Environment); err != nil {
		return ImportOutput{}, err
	}

	var addresses []importAddress
	var err error

	if input.IsApp {
		// Find the app in the definition.
		var app *appdef.App
		for i := range t.appDef.Apps {
			if t.appDef.Apps[i].Name == input.ResourceName {
				app = &t.appDef.Apps[i]
				break
			}
		}
		if app == nil {
			return ImportOutput{}, fmt.Errorf("app %q not found in app.json", input.ResourceName)
		}

		// Build import addresses based on app type and provider.
		addresses, err = buildAppImportAddresses(t.appDef.Project.Name, app, input.ResourceID)
		if err != nil {
			return ImportOutput{}, err
		}
	} else {
		// Find the resource in the definition.
		var resource *appdef.Resource
		for i := range t.appDef.Resources {
			if t.appDef.Resources[i].Name == input.ResourceName {
				resource = &t.appDef.Resources[i]
				break
			}
		}
		if resource == nil {
			return ImportOutput{}, fmt.Errorf("resource %q not found in app.json", input.ResourceName)
		}

		// Build import addresses based on resource type and provider.
		// Pass project name to build full resource names matching Terraform's naming convention.
		addresses, err = buildImportAddresses(t.appDef.Project.Name, resource, input.ResourceID)
		if err != nil {
			return ImportOutput{}, err
		}
	}

	var outputBuf strings.Builder
	t.tf.SetStdout(&outputBuf)
	t.tf.SetStderr(&outputBuf)

	var vars []tfexec.ImportOption
	for _, v := range t.env.varStrings() {
		vars = append(vars, tfexec.Var(v))
	}

	imported := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		if err = t.tf.Import(ctx, addr.Address, addr.ID, vars...); err != nil {
			return ImportOutput{
				ImportedResources: imported,
				Output:            outputBuf.String(),
			}, fmt.Errorf("importing %s: %w", addr.Address, err)
		}
		imported = append(imported, addr.Address)
	}

	return ImportOutput{
		ImportedResources: imported,
		Output:            outputBuf.String(),
	}, nil
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

func (t *Terraform) prepareVars(ctx context.Context, env env.Environment) error {
	if err := t.hasInitialised(); err != nil {
		return err
	}

	vars, err := t.tfVarsFromDefinition(ctx, env)
	if err != nil {
		return errors.Wrap(err, "generating terraform variables")
	}

	if err = t.writeTFVarsFile(vars); err != nil {
		return errors.Wrap(err, "writing tfvars file")
	}

	return nil
}

// determineImageTag determines the appropriate image tag for an app.
// Priority:
//  1. GITHUB_SHA environment variable (when running in CI)
//  2. Latest sha-* tag from GHCR via ghapi client (when running locally)
//  3. "latest" as fallback
func (t *Terraform) determineImageTag(ctx context.Context, appName string) string {
	// Check if we're in CI with GITHUB_SHA env var.
	if sha := os.Getenv("GITHUB_SHA"); sha != "" {
		return "sha-" + sha
	}

	// Try to get the latest sha tag from GHCR using the injected client.
	tag, err := t.ghClient.GetLatestSHATag(ctx, t.appDef.Project.Repo.Owner, t.appDef.Project.Repo.Name, appName)
	if err != nil {
		logger.Error("Obtaining latest SHA tag for app %Q, error: %s", appName, err.Error())

		// Fallback to latest.
		return "latest"
	}

	return tag
}

func getTerraformPath(ctx context.Context) (string, error) {
	whichCmd := executil.NewCommand("which", "terraform")
	run, err := executil.DefaultRunner().Run(ctx, whichCmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(run.Output), nil
}
