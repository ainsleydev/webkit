package cicd

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/util/executil"
	"github.com/ainsleydev/webkit/internal/util/testutil"
)

func setup(t *testing.T, fs afero.Fs, appDef *appdef.Definition) cmdtools.CommandInput {
	t.Helper()

	if !executil.Exists("action-validator") {
		t.Skip("action-validator CLI not found in PATH; skipping integration test")
	}

	return cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: appDef,
		Manifest:    manifest.NewTracker(),
	}
}

func validateWorkflow(t *testing.T, file []byte) error {
	t.Helper()

	t.Log("YAML is valid")
	if err := testutil.ValidateYAML(t, file); err != nil {
		return err
	}

	t.Log("Github Action is validated")
	if err := testutil.ValidateGithubAction(t, file); err != nil {
		return err
	}

	return nil
}
