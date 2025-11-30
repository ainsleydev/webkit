package templates

import (
	"fmt"
	"strings"
	"unicode"
)

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

// prettyConfigKey converts a snake_case configuration key to Title Case for display.
// Examples: "server_type" → "Server Type", "location" → "Location"
func prettyConfigKey(key string) string {
	// Replace underscores with spaces
	words := strings.Split(key, "_")

	// Capitalize first letter of each word
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	return strings.Join(words, " ")
}
