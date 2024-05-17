package testutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertHTML compares two HTML strings, ignoring whitespace and case sensitivity.
// It does this by:
//   - Lower-casing both strings.
//   - Removing all whitespace characters (spaces, newlines, tabs, carriage returns).
//   - Comparing the resulting strings using the `assert.Equal` function from the testing package.
func AssertHTML(t *testing.T, want, got string) {
	// Lowercase both strings
	w := strings.ToLower(want)
	f := strings.ToLower(got)

	// Remove all whitespace characters
	r := strings.NewReplacer(" ", "", "\n", "", "\t", "", "\r", "")

	// Compare the resulting strings
	assert.Equal(t, r.Replace(w), r.Replace(f))
}
