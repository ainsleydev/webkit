package jsonformat

import (
	"bytes"
	"regexp"
)

// fieldNameRegex extracts the field name from a JSON line.
var fieldNameRegex = regexp.MustCompile(`^\s*"([^"]+)"\s*:`)

// Format processes JSON output from json.MarshalIndent and applies custom
// formatting rules to make specific objects more compact.
//
// It inlines objects that contain only environment variable fields
// (source, value, path) or command fields (command, skip_ci, timeout).
func Format(data []byte) ([]byte, error) {
	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this line starts an object that should be inlined.
		if shouldInline, fieldsToInline := shouldInlineObject(lines, i); shouldInline {
			inlined, consumed := inlineObject(lines, i, fieldsToInline)
			result = append(result, inlined)
			i += consumed
		} else {
			result = append(result, line)
			i++
		}
	}

	return bytes.Join(result, []byte("\n")), nil
}

// shouldInlineObject checks if an object at the given index should be inlined.
// Returns true if it should be inlined, along with the allowed field names.
func shouldInlineObject(lines [][]byte, startIdx int) (bool, map[string]bool) {
	line := lines[startIdx]

	// Must be a line that starts an object: "key": {
	if !bytes.Contains(line, []byte(": {")) {
		return false, nil
	}

	// Look ahead to see what fields this object contains.
	if startIdx+1 >= len(lines) {
		return false, nil
	}

	// Scan the object to see what fields it has.
	fields := make(map[string]bool)
	idx := startIdx + 1
	depth := 1

	for idx < len(lines) && depth > 0 {
		currentLine := lines[idx]
		trimmed := bytes.TrimSpace(currentLine)

		// Track depth.
		if bytes.Contains(trimmed, []byte("{")) {
			depth++
		}
		if bytes.HasPrefix(trimmed, []byte("}")) {
			depth--
			if depth == 0 {
				break
			}
		}

		// Extract field name if this is a field line.
		if depth == 1 && !bytes.HasPrefix(trimmed, []byte("}")) {
			fieldName := extractFieldName(currentLine)
			if fieldName != "" {
				fields[fieldName] = true
			}
		}

		idx++
	}

	// Check if this matches an inline pattern.
	if isEnvValueObject(fields) {
		return true, envValueFields
	}
	if isCommandObject(fields) {
		return true, commandFields
	}

	return false, nil
}

var (
	// envValueFields are the allowed fields in environment value objects.
	envValueFields = map[string]bool{
		"source": true,
		"value":  true,
		"path":   true,
	}

	// commandFields are the allowed fields in command objects.
	commandFields = map[string]bool{
		"command": true,
		"skip_ci": true,
		"timeout": true,
	}
)

// isEnvValueObject checks if fields match an environment value object.
func isEnvValueObject(fields map[string]bool) bool {
	if len(fields) == 0 {
		return false
	}

	for field := range fields {
		if !envValueFields[field] {
			return false
		}
	}

	return true
}

// isCommandObject checks if fields match a command object.
func isCommandObject(fields map[string]bool) bool {
	if len(fields) == 0 {
		return false
	}

	for field := range fields {
		if !commandFields[field] {
			return false
		}
	}

	return true
}

// inlineObject takes a multi-line object and collapses it to a single line.
func inlineObject(lines [][]byte, startIdx int, allowedFields map[string]bool) ([]byte, int) {
	startLine := lines[startIdx]
	indentation := extractIndentation(startLine)
	keyPart := string(bytes.TrimSpace(startLine))

	var parts []string
	consumed := 1
	idx := startIdx + 1

	for idx < len(lines) {
		line := lines[idx]
		consumed++

		trimmed := bytes.TrimSpace(line)

		// Found the closing brace.
		if bytes.HasPrefix(trimmed, []byte("}")) {
			closeBrace := string(trimmed)
			// Build the result.
			result := indentation + keyPart
			for _, part := range parts {
				result += part + ", "
			}
			result = result[:len(result)-2] // Remove last ", "
			result += closeBrace
			return []byte(result), consumed
		}

		// Extract field part.
		fieldPart := extractFieldPart(line)
		parts = append(parts, fieldPart)

		idx++
	}

	// If we get here, something went wrong. Return original line.
	return startLine, 1
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
func extractFieldPart(line []byte) string {
	trimmed := bytes.TrimLeft(line, " \t")
	trimmed = bytes.TrimRight(trimmed, ",\r\n\t ")
	return string(trimmed)
}
