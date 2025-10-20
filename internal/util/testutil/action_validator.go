package testutil

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// ValidateYAML checks whether the bytes contain valid YAML syntax.
func ValidateYAML(t *testing.T, data []byte) error {
	t.Helper()
	var out any
	if err := yaml.Unmarshal(data, &out); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}
	return nil
}

// ValidateGithubAction validates a GitHub Actions YAML file using
// action-validator. If it conforms to GitHub actions spec,
// no error will be returned.
//
// Ref: https://github.com/mpalmer/action-validator
func ValidateGithubAction(t *testing.T, data []byte, isAction bool) error {
	t.Helper()

	// Check if action-validator is installed
	if _, err := exec.LookPath("action-validator"); err != nil {
		return errors.New("action-validator is not installed; see: https://github.com/mpalmer/action-validator")
	}

	// Determine filename based on type
	filename := "temp.yml"
	if isAction {
		filename = "action.yml"
	}

	tmpFilePath := filepath.Join(os.TempDir(), filename)
	if err := os.WriteFile(tmpFilePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	defer func(name string) {
		_ = os.Remove(name) //nolint
	}(tmpFilePath) // Ensure cleanup

	// Run the action-validator command
	cmd := exec.Command("action-validator", tmpFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("validation failed:\n%s", string(output))
	}

	return nil
}
