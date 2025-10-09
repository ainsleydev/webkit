package operations

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/config"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

const (
	webkitInfraRepo    = "https://github.com/webkit/infra.git"
	webkitInfraVersion = "v1.0.0" // TODO: Match webkit CLI version
)

// InfraPlan clones the webkit-infra repository and runs terraform plan locally
func InfraPlan(ctx context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	// 1. Prepare webkit config directory (~/.config/webkit)
	infraDir, err := prepareInfraDirectory()
	if err != nil {
		return fmt.Errorf("preparing infra directory: %w", err)
	}

	fmt.Printf("📁 Using infra directory: %s\n", infraDir)

	// 2. Clone/update webkit-infra repository
	if err := ensureInfraRepo(infraDir); err != nil {
		return fmt.Errorf("ensuring infra repo: %w", err)
	}

	// 3. Write app.json to infra directory
	gen := scaffold.New(afero.NewOsFs())
	tfVarsPath := filepath.Join(infraDir, "project.auto.tfvars.json")
	if err := gen.JSON(tfVarsPath, app); err != nil {
		return fmt.Errorf("writing terraform vars: %w", err)
	}

	fmt.Println("✓ Generated project.auto.tfvars.json")

	// 4. Configure backend (safe for local dev)
	backendPath := filepath.Join(infraDir, "backend.tf")
	if err := writeLocalBackend(backendPath, app.Project.Name); err != nil {
		return fmt.Errorf("configuring backend: %w", err)
	}

	fmt.Println("✓ Configured local backend")

	// 5. Find terraform binary
	terraformPath, err := findTerraformBinary()
	if err != nil {
		return fmt.Errorf("finding terraform: %w", err)
	}

	// 6. Initialize Terraform client
	tf, err := tfexec.NewTerraform(infraDir, terraformPath)
	if err != nil {
		return fmt.Errorf("initializing terraform: %w", err)
	}

	// Set stdout/stderr so user sees output
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	// 7. Run terraform init
	fmt.Println("\n🔧 Initializing Terraform...")
	if err := tf.Init(ctx, tfexec.Upgrade(false)); err != nil {
		return fmt.Errorf("terraform init: %w", err)
	}

	// 8. Run terraform plan
	fmt.Println("\n📋 Running Terraform Plan...")
	hasChanges, err := tf.Plan(ctx)
	if err != nil {
		return fmt.Errorf("terraform plan: %w", err)
	}

	if hasChanges {
		fmt.Println("\n⚠️  Changes detected!")
	} else {
		fmt.Println("\n✓ No changes detected")
	}

	return nil
}

// prepareInfraDirectory ensures ~/.config/webkit/infra exists
func prepareInfraDirectory() (string, error) {
	configDir, err := config.Dir()
	if err != nil {
		return "", err
	}

	infraDir := filepath.Join(configDir, "infra")
	if err := os.MkdirAll(infraDir, 0755); err != nil {
		return "", fmt.Errorf("creating infra directory: %w", err)
	}

	return infraDir, nil
}

// ensureInfraRepo clones or updates the webkit-infra repository
func ensureInfraRepo(infraDir string) error {
	gitDir := filepath.Join(infraDir, ".git")

	// Check if repo already exists
	if _, err := os.Stat(gitDir); err == nil {
		fmt.Println("📦 Updating webkit-infra repository...")
		return updateRepo(infraDir)
	}

	// Clone fresh
	fmt.Println("📦 Cloning webkit-infra repository...")
	return cloneRepo(infraDir)
}

// cloneRepo clones the webkit-infra repository
func cloneRepo(infraDir string) error {
	cmd := exec.Command("git", "clone",
		"--depth", "1",
		"--branch", webkitInfraVersion,
		webkitInfraRepo,
		infraDir,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// updateRepo pulls latest changes
func updateRepo(infraDir string) error {
	cmd := exec.Command("git", "pull", "origin", webkitInfraVersion)
	cmd.Dir = infraDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// writeLocalBackend configures Terraform to use local state for development
func writeLocalBackend(path string, projectName string) error {
	backend := fmt.Sprintf(`terraform {
  backend "local" {
    path = "terraform-%s.tfstate"
  }
}
`, projectName)

	return os.WriteFile(path, []byte(backend), 0644)
}

// findTerraformBinary locates the terraform executable
func findTerraformBinary() (string, error) {
	// Check if terraform is in PATH
	path, err := exec.LookPath("terraform")
	if err != nil {
		return "", fmt.Errorf("terraform not found in PATH. Please install terraform: https://www.terraform.io/downloads")
	}
	return path, nil
}
