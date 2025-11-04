package jsonformat

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormat_WithRealAppJSON tests the formatter with a structure similar to the real app.json.
func TestFormat_WithRealAppJSON(t *testing.T) {
	// This test simulates what happens in the webkit update command.
	// It creates a structure similar to app.json, marshals it, and applies our formatter.

	type envValue struct {
		Source string `json:"source"`
		Value  string `json:"value,omitempty"`
		Path   string `json:"path,omitempty"`
	}

	type commandSpec struct {
		Command string `json:"command,omitempty"`
		SkipCI  bool   `json:"skip_ci,omitempty"`
		Timeout string `json:"timeout,omitempty"`
	}

	type app struct {
		Name     string                     `json:"name"`
		Env      map[string]map[string]envValue `json:"env"`
		Commands map[string]commandSpec     `json:"commands"`
	}

	testApp := app{
		Name: "cms",
		Env: map[string]map[string]envValue{
			"dev": {
				"DATABASE_URI": {Source: "value", Value: "file:./cms.db"},
				"FRONTEND_URL": {Source: "value", Value: "http://localhost:5173"},
			},
			"production": {
				"DATABASE_URI": {Source: "resource", Value: "db.connection_url"},
				"FRONTEND_URL": {Source: "value", Value: "https://searchspares.com"},
			},
		},
		Commands: map[string]commandSpec{
			"build": {Command: "pnpm build"},
			"test":  {Command: "pnpm test"},
			"lint":  {Command: "pnpm lint"},
		},
	}

	// Marshal with json.MarshalIndent (what the real code does).
	data, err := json.MarshalIndent(testApp, "", "\t")
	require.NoError(t, err)

	// Apply our formatter.
	formatted, err := Format(data)
	require.NoError(t, err)

	formattedStr := string(formatted)

	// Verify environment variables are inlined.
	assert.Contains(t, formattedStr, `"DATABASE_URI": {"source": "value", "value": "file:./cms.db"}`)
	assert.Contains(t, formattedStr, `"FRONTEND_URL": {"source": "value", "value": "http://localhost:5173"}`)
	assert.Contains(t, formattedStr, `"DATABASE_URI": {"source": "resource", "value": "db.connection_url"}`)

	// Verify commands are inlined.
	assert.Contains(t, formattedStr, `"build": {"command": "pnpm build"}`)
	assert.Contains(t, formattedStr, `"test": {"command": "pnpm test"}`)
	assert.Contains(t, formattedStr, `"lint": {"command": "pnpm lint"}`)

	// Verify multi-line format is NOT present (no individual lines with just "source":).
	assert.NotContains(t, formattedStr, "\n\t\t\t\"source\": \"value\",\n")

	// Print for manual inspection if needed.
	t.Logf("Formatted output:\n%s", formattedStr)
}

// TestFormat_PreservesOtherStructures ensures we don't inline things we shouldn't.
func TestFormat_PreservesOtherStructures(t *testing.T) {
	type config struct {
		Name   string            `json:"name"`
		Nested map[string]string `json:"nested"`
	}

	testConfig := config{
		Name: "test",
		Nested: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	data, err := json.MarshalIndent(testConfig, "", "\t")
	require.NoError(t, err)

	formatted, err := Format(data)
	require.NoError(t, err)

	formattedStr := string(formatted)

	// These should still be multi-line because they don't match our inline patterns.
	assert.True(t, strings.Contains(formattedStr, "\"key1\": \"value1\""))
	assert.True(t, strings.Contains(formattedStr, "\"key2\": \"value2\""))

	// Should NOT be inlined on one line.
	assert.NotContains(t, formattedStr, `"nested": {"key1": "value1", "key2": "value2"}`)
}
