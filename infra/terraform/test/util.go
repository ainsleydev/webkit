package integration

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-json"
)

// setupTerraform configures Terraform options for integration tests.
// Set TERRAFORM_DEBUG=true to see full Terraform output.
func setupTerraform(t *testing.T, varsFixtureFile string) *terraform.Options {
	t.Helper()

	opts := &terraform.Options{
		TerraformDir: "../",
		VarFiles:     []string{"./fixtures/" + varsFixtureFile},
		PlanFilePath: t.TempDir() + "/tfplan.out",
	}

	if os.Getenv("TERRAFORM_DEBUG") != "true" {
		opts.Logger = logger.Discard
		opts.NoColor = true
	}

	return opts
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
