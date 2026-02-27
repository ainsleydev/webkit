package appdef

import (
	"fmt"
	"path/filepath"
)

type (
	// UtilityCI defines the CI configuration for a utility.
	// When present, a CI job will be generated for this utility.
	// Omit CI entirely to create a workspace-only utility with no CI job.
	UtilityCI struct {
		Trigger  string `json:"trigger" validate:"required,oneof=pull_request cron" description:"Event that triggers this CI job (pull_request or cron)"`
		Schedule string `json:"schedule,omitempty" validate:"required_if=Trigger cron" description:"Cron expression for scheduled triggers (e.g. '0 0 * * *')"`
		RunsOn   string `json:"runs_on,omitempty" description:"GitHub Actions runner (defaults to ubuntu-latest)"`
	}

	// Utility represents a non-deployed workspace member within the webkit project.
	// Utilities are included in the pnpm workspace (if JS) and optionally run on CI,
	// but are never deployed. Examples include E2E tests, shared constants,
	// benchmark suites, and CLI tools.
	Utility struct {
		Name        string     `json:"name" validate:"required,lowercase,alphanumdash" description:"Unique identifier for the utility (lowercase, hyphenated)"`
		Title       string     `json:"title" validate:"required" description:"Human-readable utility name for display purposes"`
		Description string     `json:"description,omitempty" validate:"omitempty,max=200" description:"Brief description of the utility's purpose and functionality"`
		Path        string     `json:"path" validate:"required" description:"Relative file path to the utility's source code directory"`
		Language    string     `json:"language" validate:"required,oneof=go js" description:"Toolchain language for CI setup and workspace inclusion (go or js)"`
		CI          *UtilityCI `json:"ci,omitempty" description:"CI configuration. Omit to create a workspace-only utility with no CI job"`
		Toolset
	}
)

// HasCI returns whether this utility has CI configuration and should
// generate a CI job.
func (u *Utility) HasCI() bool {
	return u.CI != nil
}

// ShouldUseNPM returns whether this utility should be included in
// the pnpm workspace. JS utilities are included, Go utilities are not.
func (u *Utility) ShouldUseNPM() bool {
	return u.Language == "js"
}

// applyDefaults sets default values for the utility.
func (u *Utility) applyDefaults() error {
	u.Toolset.initDefaults()
	if u.Path != "" {
		u.Path = filepath.Clean(u.Path)
	}
	if u.CI != nil {
		if u.CI.RunsOn == "" {
			u.CI.RunsOn = "ubuntu-latest"
		}
		if u.CI.Trigger == "cron" && u.CI.Schedule == "" {
			return fmt.Errorf("utility %q has cron trigger but no schedule", u.Name)
		}
	}
	return nil
}
