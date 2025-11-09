package stringutil

import (
	"regexp"
	"strings"

	"github.com/yosssi/gohtml"
)

// IsNotEmpty reports whether the string pointer is non-nil and non-empty.
func IsNotEmpty(in *string) bool {
	return in != nil && *in != ""
}

// IsEmpty reports whether the string pointer is nil or empty.
func IsEmpty(in *string) bool {
	return in == nil || *in == ""
}

var space = regexp.MustCompile(`\s+`)

// RemoveDuplicateWhitespace collapses consecutive whitespace into single spaces
// and trims leading/trailing whitespace.
func RemoveDuplicateWhitespace(s string) string {
	return strings.TrimSpace(space.ReplaceAllString(s, " "))
}

// FormatHTML formats HTML content with proper indentation and structure.
func FormatHTML(s string) string {
	repl := strings.NewReplacer(
		" >", ">",
	)
	return gohtml.Format(repl.Replace(s))
}
