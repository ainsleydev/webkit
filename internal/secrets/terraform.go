package secrets

import (
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/pkg/env"
)

// TransformOutputs converts an OutputResult from Terraform into a
// TerraformOutputProvider that can be used for secret resolution.
//
// This function extracts resource outputs and creates OutputKeys for
// each environment/resource/output combination.
func TransformOutputs(result infra.OutputResult, environment env.Environment) TerraformOutputProvider {
	provider := make(TerraformOutputProvider)

	for resourceName, outputs := range result.Resources {
		for outputName, value := range outputs {
			key := OutputKey{
				Environment:  environment,
				ResourceName: resourceName,
				OutputName:   outputName,
			}
			provider[key] = value
		}
	}

	return provider
}
