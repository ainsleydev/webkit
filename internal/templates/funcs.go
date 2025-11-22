package templates

import "fmt"

// githubExpression returns a GitHub Actions expression.
// Use for context expressions like github.sha, steps.*.outputs.*, etc.
func githubExpression(in string) string {
	return fmt.Sprintf("${{ %s }}", in)
}

// githubVariable wraps a variable name in GitHub Actions syntax.
// Use for repository or organization variables.
func githubVariable(name string) string {
	return fmt.Sprintf("${{ vars.%s }}", name)
}

// githubSecret wraps a secret name in GitHub Actions syntax.
func githubSecret(name string) string {
	return fmt.Sprintf("${{ secrets.%s }}", name)
}

// githubInput wraps an input name in GitHub Actions syntax.
func githubInput(name string) string {
	return fmt.Sprintf("${{ inputs.%s }}", name)
}

// githubEnv wraps an env name in GitHub Actions syntax.
func githubEnv(name string) string {
	return fmt.Sprintf("${{ env.%s }}", name)
}
