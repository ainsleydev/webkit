package jsonformat

import (
	"bytes"
	"regexp"
	"strings"
)

// Format processes JSON output from json.MarshalIndent and applies custom
// formatting rules to make specific objects more compact.
//
// Specifically, it collapses:
// - Environment variable values (source, value, path fields)
// - Command specifications (command, skip_ci, timeout fields)
//
// This improves readability by reducing vertical space for simple objects.
func Format(data []byte) ([]byte, error) {
	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this line starts an object that should be inlined.
		if shouldInline, pattern := isInlineCandidate(line, lines, i); shouldInline {
			inlined, consumed := inlineObject(lines, i, pattern)
			result = append(result, inlined)
			i += consumed
		} else {
			result = append(result, line)
			i++
		}
	}

	return bytes.Join(result, []byte("\n")), nil
}

// inlinePattern defines the structure of objects that should be collapsed.
type inlinePattern struct {
	// Fields that are allowed in the inline object.
	allowedFields map[string]bool
	// Whether this is a simple object (only 2 fields typically).
	isSimple bool
}

var (
	// envValuePattern matches environment variable values.
	envValuePattern = inlinePattern{
		allowedFields: map[string]bool{
			"source": true,
			"value":  true,
			"path":   true,
		},
		isSimple: true,
	}

	// commandPattern matches command specifications.
	commandPattern = inlinePattern{
		allowedFields: map[string]bool{
			"command": true,
			"skip_ci": true,
			"timeout": true,
		},
		isSimple: false,
	}
)

// fieldNameRegex extracts the field name from a JSON line.
// Matches: "field_name":
var fieldNameRegex = regexp.MustCompile(`^\s*"([^"]+)"\s*:`)

// isInlineCandidate checks if a line starts an object that should be inlined.
func isInlineCandidate(line []byte, lines [][]byte, idx int) (bool, *inlinePattern) {
	// Must be a line that starts an object: "key": {
	if !bytes.Contains(line, []byte(": {")) {
		return false, nil
	}

	// Look ahead to see what fields this object contains.
	if idx+1 >= len(lines) {
		return false, nil
	}

	nextLine := lines[idx+1]
	fieldName := extractFieldName(nextLine)

	// Check if the next line contains a field from our inline patterns.
	if envValuePattern.allowedFields[fieldName] {
		return isValidInlineObject(lines, idx, &envValuePattern), &envValuePattern
	}

	if commandPattern.allowedFields[fieldName] {
		return isValidInlineObject(lines, idx, &commandPattern), &commandPattern
	}

	return false, nil
}

// isValidInlineObject verifies that all fields in the object match the pattern.
func isValidInlineObject(lines [][]byte, startIdx int, pattern *inlinePattern) bool {
	// Scan ahead to ensure all fields match the pattern.
	idx := startIdx + 1
	for idx < len(lines) {
		line := lines[idx]

		// Found the closing brace.
		if isClosingBrace(line) {
			return true
		}

		// Check if this field is allowed.
		fieldName := extractFieldName(line)
		if fieldName == "" || !pattern.allowedFields[fieldName] {
			return false
		}

		idx++
	}

	return false
}

// inlineObject takes a multi-line object and collapses it to a single line.
// Returns the inlined bytes and the number of lines consumed.
func inlineObject(lines [][]byte, startIdx int, pattern *inlinePattern) ([]byte, int) {
	startLine := lines[startIdx]
	indentation := extractIndentation(startLine)

	// Extract the key and opening brace without indentation.
	keyPart := strings.TrimSpace(string(startLine))

	var parts []string
	parts = append(parts, keyPart) // Start with "key": {

	consumed := 1
	idx := startIdx + 1

	for idx < len(lines) {
		line := lines[idx]
		consumed++

		// Found the closing brace.
		if isClosingBrace(line) {
			// Add the closing brace.
			closeBrace := extractClosingBrace(line)
			parts = append(parts, closeBrace)
			break
		}

		// Extract the field part.
		fieldPart := extractFieldPart(line)
		parts = append(parts, fieldPart)

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

// extractFieldName extracts the field name from a JSON line.
func extractFieldName(line []byte) string {
	matches := fieldNameRegex.FindSubmatch(line)
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

// extractFieldPart extracts the field definition without indentation.
// Example: `    "source": "value",` -> `"source": "value"`
func extractFieldPart(line []byte) string {
	trimmed := bytes.TrimLeft(line, " \t")
	// Remove trailing comma and whitespace.
	trimmed = bytes.TrimRight(trimmed, ",\r\n\t ")
	return string(trimmed)
}

// isClosingBrace checks if a line contains only a closing brace with optional comma.
func isClosingBrace(line []byte) bool {
	trimmed := bytes.TrimSpace(line)
	return bytes.HasPrefix(trimmed, []byte("}"))
}

// extractClosingBrace returns the closing brace with optional trailing comma.
func extractClosingBrace(line []byte) string {
	trimmed := bytes.TrimSpace(line)
	return string(trimmed)
}
