package stringutil

import (
	"regexp"
	"strings"
)

// IsNotEmpty checks if a string pointer is empty.
func IsNotEmpty(in *string) bool {
	return in != nil && *in != ""
}

// IsEmpty checks if a string pointer is empty.
func IsEmpty(in *string) bool {
	return in == nil || *in == ""
}

var space = regexp.MustCompile(`\s+`)

// RemoveDuplicateWhitespace removes duplicate whitespace from a string.
// Including tab & new line characters.
func RemoveDuplicateWhitespace(s string) string {
	return strings.TrimSpace(space.ReplaceAllString(s, " "))
}
