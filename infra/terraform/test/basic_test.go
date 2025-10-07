package integration

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformBasicValidation(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform",
		NoColor:      true,
	}

	output, err := terraform.ValidateE(t, terraformOptions)
	assert.NoError(t, err, "Terraform configuration should be valid")
	assert.Contains(t, output, "Success")
}

func TestTerraformPlanMinimal(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform",
		VarFiles:     []string{"../test/fixtures/minimal.tfvars"},
		PlanFilePath: t.TempDir() + "/tfplan.out",
	}

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err, "Terraform plan should be successful")

	// Example assertion: 5 resources to add
	assert.Equal(t, 5, len(plan.ResourcePlannedValuesMap))

	for k, _ := range plan.ResourcePlannedValuesMap {
		fmt.Println("-------")
		fmt.Println(k)
	}

}
