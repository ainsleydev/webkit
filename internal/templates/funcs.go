package templates

import "fmt"

// githubSecret wraps a secret name in GitHub Actions syntax.
func githubSecret(name string) string {
	return fmt.Sprintf("${{ secrets.%s }}", name)
}

// githubSecret wraps a input name in GitHub Actions syntax.
func githubInput(name string) string {
	return fmt.Sprintf("${{ inputs.%s }}", name)
}
