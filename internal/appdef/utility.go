package appdef

import (
	"fmt"
	"path/filepath"
	"regexp"
)

type (
	// UtilityCI defines the CI configuration for a utility.
	// When present, a CI job will be generated for this utility.
	// Omit CI entirely to create a workspace-only utility with no CI job.
	UtilityCI struct {
		Trigger  string `json:"trigger" required:"true" validate:"required,oneof=pull_request cron" enum:"pull_request,cron" description:"Event that triggers this CI job (pull_request or cron)"`
		Schedule string `json:"schedule,omitempty" validate:"required_if=Trigger cron" description:"Cron expression for scheduled triggers (e.g. '0 0 * * *')"`
		RunsOn   string `json:"runs_on,omitempty" description:"GitHub Actions runner (defaults to ubuntu-latest)"`
	}
	// Utility represents a non-deployed workspace member within the webkit project.
	// Utilities are included in the pnpm workspace (if JS) and optionally run on CI,
	// but are never deployed. Examples include E2E tests, shared constants,
	// benchmark suites, and CLI tools.
	Utility struct {
		Name        string     `json:"name" required:"true" validate:"required,lowercase,alphanumdash" description:"Unique identifier for the utility (lowercase, hyphenated)"`
		Title       string     `json:"title" required:"true" validate:"required" description:"Human-readable utility name for display purposes"`
		Description string     `json:"description,omitempty" validate:"omitempty,max=200" description:"Brief description of the utility's purpose and functionality"`
		Path        string     `json:"path" required:"true" validate:"required" description:"Relative file path to the utility's source code directory"`
		Language    string     `json:"language" required:"true" validate:"required,oneof=go js" enum:"go,js" description:"Toolchain language for CI setup and workspace inclusion (go or js)"`
		CI          *UtilityCI `json:"ci,omitempty" description:"CI configuration. Omit to create a workspace-only utility with no CI job"`
		Toolset
	}
)

// cronRegexp validates a 5-field cron expression (minute hour dom month dow).
// Each field allows digits, *, commas, hyphens, and slashes.
var cronRegexp = regexp.MustCompile(`^(\*|[0-9,\-/]+)\s+(\*|[0-9,\-/]+)\s+(\*|[0-9,\-/]+)\s+(\*|[0-9,\-/]+)\s+(\*|[0-9,\-/]+)$`)

// nameAndPath returns the utility's name and path, satisfying the pathItem interface.
func (u Utility) nameAndPath() (string, string) { return u.Name, u.Path }

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
		if u.CI.Trigger == "cron" {
			if u.CI.Schedule == "" {
				return fmt.Errorf("utility %q has cron trigger but no schedule", u.Name)
			}
			if !cronRegexp.MatchString(u.CI.Schedule) {
				return fmt.Errorf("utility %q has invalid cron expression %q: expected 5 fields (minute hour dom month dow), e.g. '0 2 * * 1'", u.Name, u.CI.Schedule)
			}
		}
	}
	return nil
}
