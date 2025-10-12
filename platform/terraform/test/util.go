package integration

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/require"
)

// setupTerraform configures Terraform options for integration tests.
// Set TERRAFORM_DEBUG=true to see full Terraform output.
func setupTerraform(t *testing.T, varsFixtureFile string) (*terraform.Options, func()) {
	t.Helper()

	// Create a temp directory for plan files
	tempDir, err := os.MkdirTemp("", "tfplan")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	opts := &terraform.Options{
		TerraformDir: "../",
		VarFiles:     []string{"./fixtures/" + varsFixtureFile},
		PlanFilePath: tempDir + "/tfplan.out",
	}

	//if os.Getenv("TERRAFORM_DEBUG") != "true" {
	//	opts.Logger = logger.Discard
	//	opts.NoColor = true
	//}

	return opts, func() {
		require.NoError(t, os.RemoveAll(tempDir))
	}
}

// findResource locates a Terraform resource in the plan output by searching
// for a substring match in the resource address.
// Returns an error if the resource cannot be found.
func findResource(
	resourceType string,
	plannedResources map[string]*tfjson.StateResource,
) (*tfjson.StateResource, error) {
	for address, resource := range plannedResources {
		if strings.Contains(address, resourceType) {
			return resource, nil
		}
	}
	return nil, errors.New("resource not found: " + resourceType)
}
