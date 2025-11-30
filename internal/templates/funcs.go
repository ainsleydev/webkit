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
// It handles edge cases like empty strings, consecutive underscores, and leading/trailing underscores.
//
// Examples:
//   - "server_type" → "Server Type"
//   - "location" → "Location"
//   - "mount_point" → "Mount Point"
//   - "" → ""
func prettyConfigKey(key string) string {
	// Handle empty string edge case
	if key == "" {
		return ""
	}

	// Split the key by underscores to get individual words
	words := strings.Split(key, "_")

	// Capitalize the first letter of each word
	for i, word := range words {
		// Skip empty segments (from consecutive/leading/trailing underscores)
		if len(word) == 0 {
			continue
		}

		// Convert to runes to handle Unicode properly
		runes := []rune(word)
		// Capitalize first character
		runes[0] = unicode.ToUpper(runes[0])
		// Store the capitalized word back
		words[i] = string(runes)
	}

	return strings.Join(words, " ")
}
