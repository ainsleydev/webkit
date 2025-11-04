package jsonformat

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
)

// inlineParents contains JSON field names whose child objects should be inlined.
// Built at init time by scanning struct tags.
var inlineParents map[string]bool

func init() {
	inlineParents = make(map[string]bool)
}

// RegisterType scans a struct type for fields tagged with inline:"true".
func RegisterType(t reflect.Type) {
	scanType(t)
}

// scanType recursively scans a type for inline tags.
func scanType(t reflect.Type) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			continue
		}

		// If this field has inline:"true", its children should be inlined.
		if field.Tag.Get("inline") == "true" {
			inlineParents[jsonName] = true
		}

		// Recurse into struct fields.
		fieldType := field.Type
		for fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Map {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			scanType(fieldType)
		}
	}
}

var fieldNameRegex = regexp.MustCompile(`^\s*"([^"]+)"\s*:`)

// Format processes JSON output from json.MarshalIndent and applies custom
// formatting rules to make specific objects more compact.
func Format(data []byte) ([]byte, error) {
	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte

	// Track which inline parent we're currently inside.
	var currentParents []string

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this line defines an inline parent (e.g., "dev": {).
		if parent := extractParentIfInlineCandidate(line); parent != "" {
			currentParents = append(currentParents, parent)
			result = append(result, line)
			i++
			continue
		}

		// Check if this line closes a brace, pop parent.
		if isClosingBrace(line) {
			if len(currentParents) > 0 {
				currentParents = currentParents[:len(currentParents)-1]
			}
			result = append(result, line)
			i++
			continue
		}

		// Check if we should inline this object.
		shouldInline := false

		// Check if we're inside an inline parent.
		if len(currentParents) > 0 {
			parent := currentParents[len(currentParents)-1]
			if inlineParents[parent] {
				// We're inside an inline parent, inline this object.
				shouldInline = true
			}
		}

		// Fallback: check if this object matches env/command pattern.
		if !shouldInline {
			shouldInline, _ = shouldInlineObject(lines, i)
		}

		if shouldInline {
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

// extractParentIfInlineCandidate checks if this line starts an inline parent.
func extractParentIfInlineCandidate(line []byte) string {
	if !bytes.Contains(line, []byte(": {")) {
		return ""
	}

	key := extractFieldName(line)
	if key != "" && inlineParents[key] {
		return key
	}

	return ""
}

// isClosingBrace checks if a line is just a closing brace.
func isClosingBrace(line []byte) bool {
	trimmed := bytes.TrimSpace(line)
	return bytes.HasPrefix(trimmed, []byte("}"))
}

// shouldInlineObject checks if an object matches an inline pattern.
func shouldInlineObject(lines [][]byte, startIdx int) (bool, map[string]bool) {
	line := lines[startIdx]

	if !bytes.Contains(line, []byte(": {")) {
		return false, nil
	}

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

		if bytes.Contains(trimmed, []byte("{")) {
			depth++
		}
		if bytes.HasPrefix(trimmed, []byte("}")) {
			depth--
			if depth == 0 {
				break
			}
		}

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
	envValueFields = map[string]bool{
		"source": true,
		"value":  true,
		"path":   true,
	}

	commandFields = map[string]bool{
		"command": true,
		"skip_ci": true,
		"timeout": true,
	}
)

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

func inlineObject(lines [][]byte, startIdx int) ([]byte, int) {
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

		if bytes.HasPrefix(trimmed, []byte("}")) {
			closeBrace := string(trimmed)
			result := indentation + keyPart
			for _, part := range parts {
				result += part + ", "
			}
			if len(parts) > 0 {
				result = result[:len(result)-2] // Remove last ", "
			}
			result += closeBrace
			return []byte(result), consumed
		}

		fieldPart := extractFieldPart(line)
		parts = append(parts, fieldPart)

		idx++
	}

	return startLine, 1
}

func extractIndentation(line []byte) string {
	for i, b := range line {
		if b != ' ' && b != '\t' {
			return string(line[:i])
		}
	}
	return string(line)
}

func extractFieldName(line []byte) string {
	matches := fieldNameRegex.FindSubmatch(line)
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

func extractFieldPart(line []byte) string {
	trimmed := bytes.TrimLeft(line, " \t")
	trimmed = bytes.TrimRight(trimmed, ",\r\n\t ")
	return string(trimmed)
}
