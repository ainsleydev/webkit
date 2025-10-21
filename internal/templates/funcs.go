package templates

import "fmt"

// githubVariable returns a simple variable.
func githubVariable(in string) string {
	return fmt.Sprintf("${{ %s }}", in)
}

// githubSecret wraps a secret name in GitHub Actions syntax.
func githubSecret(name string) string {
	return fmt.Sprintf("${{ secrets.%s }}", name)
}

// githubSecret wraps an input name in GitHub Actions syntax.
func githubInput(name string) string {
	return fmt.Sprintf("${{ inputs.%s }}", name)
}

// githubEnv wraps an env name in GitHub Actions syntax.
func githubEnv(name string) string {
	return fmt.Sprintf("${{ env.%s }}", name)
}
