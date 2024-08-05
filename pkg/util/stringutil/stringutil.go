package stringutil

import (
	"regexp"
	"strings"

	"github.com/yosssi/gohtml"
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

// FormatHTML parses a input HTML string, formats it and returns the result.
func FormatHTML(s string) string {
	repl := strings.NewReplacer(
		" >", ">",
	)
	return gohtml.Format(repl.Replace(s))
}
