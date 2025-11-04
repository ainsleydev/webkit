package jsonformat

import (
	"bytes"
	"reflect"
	"strings"
)

// inlineParents contains JSON field names whose child objects should be inlined.
// This is built at init time by scanning struct tags.
var inlineParents map[string]bool

// init initialises the inline parents map.
func init() {
	inlineParents = make(map[string]bool)
}

// RegisterType scans a struct type for fields tagged with inline:"true"
// and registers their JSON field names for inline formatting.
func RegisterType(t reflect.Type) {
	scanType(t)
}

// scanType recursively scans a type for inline tags.
func scanType(t reflect.Type) {
	// Dereference pointers.
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields.
		if !field.IsExported() {
			continue
		}

		// Extract JSON tag.
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse JSON field name (before comma).
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			continue
		}

		// Check for inline tag.
		if field.Tag.Get("inline") == "true" {
			inlineParents[jsonName] = true
		}

		// Recursively scan the field type.
		fieldType := field.Type
		for fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Slice {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			scanType(fieldType)
		}
	}
}

// Format processes JSON output from json.MarshalIndent and applies custom
// formatting rules to make specific objects more compact.
//
// Objects that are direct children of fields tagged with inline:"true" will be
// collapsed to single lines.
func Format(data []byte) ([]byte, error) {
	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this line starts an object that should be inlined.
		if shouldInline := isInlineCandidate(line); shouldInline {
			inlined, consumed := inlineObject(lines, i)
			result = append(result, inlined)
			i += consumed
		} else {
			result = append(result, line)
			i++
		}
	}

	return bytes.Join(result, []byte("\n")), nil
}

// isInlineCandidate checks if a line starts an object whose parent is tagged for inlining.
func isInlineCandidate(line []byte) bool {
	// Must be a line that starts an object: "key": {
	if !bytes.Contains(line, []byte(": {")) {
		return false
	}

	// Extract the key name.
	key := extractParentKey(line)
	if key == "" {
		return false
	}

	// Check if this key's children should be inlined.
	return inlineParents[key]
}

// extractParentKey extracts the parent field name from a line like: "dev": {
func extractParentKey(line []byte) string {
	// Find the opening quote.
	start := bytes.IndexByte(line, '"')
	if start == -1 {
		return ""
	}

	// Find the closing quote.
	end := bytes.IndexByte(line[start+1:], '"')
	if end == -1 {
		return ""
	}

	return string(line[start+1 : start+1+end])
}

// inlineObject takes a multi-line object and collapses it to a single line.
// Returns the inlined bytes and the number of lines consumed.
func inlineObject(lines [][]byte, startIdx int) ([]byte, int) {
	startLine := lines[startIdx]
	indentation := extractIndentation(startLine)

	// Extract the key and opening brace without indentation.
	keyPart := strings.TrimSpace(string(startLine))

	var parts []string
	parts = append(parts, keyPart) // Start with "key": {

	consumed := 1
	idx := startIdx + 1
	depth := 1

	for idx < len(lines) && depth > 0 {
		line := lines[idx]
		consumed++

		trimmed := bytes.TrimSpace(line)

		// Track brace depth to handle nested objects.
		if bytes.Contains(trimmed, []byte("{")) {
			depth++
		}
		if bytes.HasPrefix(trimmed, []byte("}")) {
			depth--
			if depth == 0 {
				// Found the closing brace for this object.
				closeBrace := extractClosingBrace(line)
				parts = append(parts, closeBrace)
				break
			}
		}

		// Extract the field part (if not a closing brace).
		if !bytes.HasPrefix(trimmed, []byte("}")) {
			fieldPart := extractFieldPart(line)
			parts = append(parts, fieldPart)
		}

		idx++
	}

	// Join with proper spacing: "key": {"field": "value", "field2": "value2"}
	result := indentation + parts[0]
	for i := 1; i < len(parts)-1; i++ {
		result += parts[i] + ", "
	}
	if len(parts) > 1 {
		// Remove trailing comma and space, add closing brace.
		result = strings.TrimSuffix(result, ", ")
		result += parts[len(parts)-1]
	}

	return []byte(result), consumed
}

// extractIndentation returns the leading whitespace from a line.
func extractIndentation(line []byte) string {
	for i, b := range line {
		if b != ' ' && b != '\t' {
			return string(line[:i])
		}
	}
	return string(line)
}

// extractFieldPart extracts the field definition without indentation.
// Example: `    "source": "value",` -> `"source": "value"`
func extractFieldPart(line []byte) string {
	trimmed := bytes.TrimLeft(line, " \t")
	// Remove trailing comma and whitespace.
	trimmed = bytes.TrimRight(trimmed, ",\r\n\t ")
	return string(trimmed)
}

// extractClosingBrace returns the closing brace with optional trailing comma.
func extractClosingBrace(line []byte) string {
	trimmed := bytes.TrimSpace(line)
	return string(trimmed)
}
