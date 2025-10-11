package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/ainsleydev/webkit/infra/terraform"
	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/config"
	"github.com/ainsleydev/webkit/internal/executil"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/git"
)

const (
	webkitInfraRepo = "https://github.com/ainsleydev/webkit-infra.git"
	webkitInfraRef  = "main" // or version tag like "v1.2.3"
)

func Test(ctx context.Context, input cmdtools.CommandInput) error {
	tmpDir, err := os.MkdirTemp("", "webkit-tf")
	if err != nil {
		return err
	}
	defer func(path string) {
		if err = os.RemoveAll(path); err != nil {
			slog.ErrorContext(ctx, "Failed to remove temp dir", slog.String("path", path))
		}
	}(tmpDir)

	fmt.Println("Temp dir:", tmpDir)

	// Copy embedded templates to the working directory
	err = fsext.CopyAllEmbed(tfembed.Templates, tmpDir)
	if err != nil {
		return err
	}

	whichCmd := executil.NewCommand("which", "terraform")
	run, err := executil.DefaultRunner().Run(ctx, whichCmd)
	if err != nil {
		return err
	}
	tfPath := strings.TrimSpace(run.Output)

	// Create terraform executor
	tfDir := filepath.Join(tmpDir, "base")
	tf, err := tfexec.NewTerraform(tfDir, tfPath)
	if err != nil {
		return fmt.Errorf("creating terraform executor: %w", err)
	}

	// Resolve env vars

	fmt.Println("Initializing Terraform...")
	if err = tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return fmt.Errorf("terraform init: %w", err)
	}

	fmt.Println("Making Plan....")
	_, err = tf.Plan(ctx,
		tfexec.Var("project_name=hey!"),
		tfexec.Var("environment=production!"),
	)
	if err != nil {
		return fmt.Errorf("executing plan: %w", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(state)

	return nil

	//fmt.Println(run.Output)
	//os.Exit(1)
	//// Run Terraform commands
	//tf, err := tfexec.NewTerraform(tmp, "/usr/local/bin/terraform")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//ctx := context.Background()
	//if err := tf.Init(ctx); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := tf.Plan(ctx); err != nil {
	//	log.Fatal(err)
	//}

}

// InfraPlan runs terraform plan using the webkit-infra repository
//
// This function:
// 1. Clones or updates the webkit-infra repository
// 2. Writes app.json to the terraform working directory
// 3. Initializes terraform with the appropriate backend
// 4. Runs terraform plan to preview infrastructure changes
//
// Note: app.json is passed to terraform unmodified. If terraform cannot
// handle the app.json structure, this will be noted but not block execution.
func InfraPlan(ctx context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	// 1. Determine where to clone webkit-infra repo
	configDir, err := config.Dir()
	if err != nil {
		return fmt.Errorf("getting config dir: %w", err)
	}

	infraPath := filepath.Join(configDir, "webkit-infra")

	// 2. Clone or update webkit-infra repository using the git package
	if err := ensureWebKitInfra(ctx, infraPath); err != nil {
		return err
	}

	// 3. Navigate to infra subdirectory (where Terraform root is)
	terraformDir := filepath.Join(infraPath, "infra")

	// 4. Write app.json to the terraform directory
	// This passes the app definition to terraform with no modifications
	if err := writeAppJSON(terraformDir, app); err != nil {
		return fmt.Errorf("writing app.json: %w", err)
	}

	// 5. Initialize Terraform
	tf, err := initTerraform(ctx, terraformDir)
	if err != nil {
		return fmt.Errorf("initializing terraform: %w", err)
	}

	// 6. Run terraform plan
	fmt.Println("Running terraform plan...")
	if err := runTerraformPlan(ctx, tf); err != nil {
		return fmt.Errorf("terraform plan failed: %w", err)
	}

	fmt.Println("Terraform plan completed successfully!")
	return nil
}

// ensureWebKitInfra clones or updates the webkit-infra repository
func ensureWebKitInfra(ctx context.Context, infraPath string) error {
	// Create git client with default runner
	gitClient, err := git.New(executil.DefaultRunner())
	if err != nil {
		return fmt.Errorf("creating git client: %w", err)
	}

	cfg := git.CloneConfig{
		URL:       webkitInfraRepo,
		LocalPath: infraPath,
		Ref:       webkitInfraRef,
		Depth:     1, // Shallow clone for faster downloads
	}

	fmt.Printf("Ensuring webkit-infra repository at %s...\n", infraPath)
	if err := gitClient.CloneOrUpdate(ctx, cfg); err != nil {
		return fmt.Errorf("cloning/updating webkit-infra: %w", err)
	}

	return nil
}

// writeAppJSON writes the app definition to app.json in the terraform directory
// The app.json is passed to terraform with no modifications
func writeAppJSON(terraformDir string, app *appdef.Definition) error {
	appJSONPath := filepath.Join(terraformDir, "app.json")

	// Marshal the app definition to JSON with indentation for readability
	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling app definition: %w", err)
	}

	// Write app.json to the terraform directory
	if err := os.WriteFile(appJSONPath, data, 0644); err != nil {
		return fmt.Errorf("writing app.json file: %w", err)
	}

	fmt.Printf("Wrote app.json to %s\n", appJSONPath)

	// NOTE: If terraform cannot handle app.json in this format, we note it here
	// but do not block execution as per requirements
	fmt.Println("Note: app.json passed to terraform unmodified")
	fmt.Println("If terraform cannot parse this format, terraform init/plan will fail")

	return nil
}

// initTerraform initializes terraform in the given directory
func initTerraform(ctx context.Context, workingDir string) (*tfexec.Terraform, error) {
	// Find terraform executable in PATH
	terraformPath, err := findTerraformExecutable()
	if err != nil {
		return nil, err
	}

	// Create terraform executor
	tf, err := tfexec.NewTerraform(workingDir, terraformPath)
	if err != nil {
		return nil, fmt.Errorf("creating terraform executor: %w", err)
	}

	// Initialize terraform
	fmt.Println("Initializing terraform...")
	if err := tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return nil, fmt.Errorf("terraform init: %w", err)
	}

	return tf, nil
}

// runTerraformPlan executes terraform plan
func runTerraformPlan(ctx context.Context, tf *tfexec.Terraform) error {
	// Run terraform plan
	// The plan output will be printed to stdout automatically
	_, err := tf.Plan(ctx)
	if err != nil {
		return fmt.Errorf("executing plan: %w", err)
	}

	return nil
}

// findTerraformExecutable locates the terraform binary in PATH
func findTerraformExecutable() (string, error) {
	// Try to find terraform in PATH
	path, err := exec.LookPath("terraform")
	if err != nil {
		return "", fmt.Errorf("terraform not found in PATH: %w", err)
	}
	return path, nil
}
