package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/util/executil"
)

func main() {
	console := printer.New(os.Stdout)

	// Main menu
	console.Printf("Tag Management\n\n")
	console.Printf("  1) Create new tag\n")
	console.Printf("  2) Delete existing tag\n")
	console.Printf("  3) Exit\n\n")

	reader := bufio.NewReader(os.Stdin)
	console.Printf("Enter choice [1-3]: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		console.Error(fmt.Sprintf("Error reading input: %s", err))
		os.Exit(1)
	}

	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		createTag(console, reader)
	case "2":
		deleteTag(console, reader)
	case "3":
		console.Printf("Goodbye!\n")
		return
	default:
		console.Error("Invalid choice")
		os.Exit(1)
	}
}

func createTag(console *printer.Console, reader *bufio.Reader) {
	runner := executil.DefaultRunner()
	ctx := context.Background()

	latestTag := getLatestTag(runner, ctx)

	console.Printf("\nCurrent version: %s\n\n", latestTag)

	// Parse current version
	current, err := semver.NewVersion(latestTag)
	if err != nil {
		console.Error(fmt.Sprintf("Error parsing version: %s", err))
		os.Exit(1)
	}

	// Calculate next versions
	nextPatch := current.IncPatch()
	nextMinor := current.IncMinor()
	nextMajor := current.IncMajor()

	// Display menu
	console.Printf("What type of release do you want to create?\n\n")
	console.Printf("  1) Patch   - v%s  (bug fixes)\n", nextPatch.String())
	console.Printf("  2) Minor   - v%s  (new features, backwards compatible)\n", nextMinor.String())
	console.Printf("  3) Major   - v%s  (breaking changes)\n", nextMajor.String())
	console.Printf("  4) Cancel\n\n")

	// Get user input
	console.Printf("Enter choice [1-4]: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		console.Error(fmt.Sprintf("Error reading input: %s", err))
		os.Exit(1)
	}

	choice = strings.TrimSpace(choice)

	var newVersion *semver.Version
	switch choice {
	case "1":
		newVersion = &nextPatch
	case "2":
		newVersion = &nextMinor
	case "3":
		newVersion = &nextMajor
	case "4":
		console.Printf("Cancelled.\n")
		return
	default:
		console.Error("Invalid choice")
		os.Exit(1)
	}

	newTag := "v" + newVersion.String()

	// Confirm
	console.Printf("\nCreate tag %s? (y/n): ", newTag)
	confirm, err := reader.ReadString('\n')
	if err != nil {
		console.Error(fmt.Sprintf("Error reading confirmation: %s", err))
		os.Exit(1)
	}

	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" && confirm != "yes" {
		console.Printf("Cancelled.\n")
		return
	}

	// Run go generate to ensure all generated files are up to date
	console.Printf("\nRunning go generate...\n")
	if err := runGoGenerate(runner, ctx); err != nil {
		console.Error(fmt.Sprintf("Error running go generate: %s", err))
		os.Exit(1)
	}
	console.Success("Generated files updated")

	// Generate version file with new tag
	console.Printf("Generating version file...\n")
	if err := generateVersionFile(console, newTag); err != nil {
		console.Error(fmt.Sprintf("Error generating version file: %s", err))
		os.Exit(1)
	}
	console.Success("Version file generated")

	// Commit all changes (generated files + version file)
	console.Printf("Committing changes...\n")
	if err := gitCommitAllChanges(runner, ctx, newTag); err != nil {
		console.Error(fmt.Sprintf("Error committing changes: %s", err))
		os.Exit(1)
	}
	console.Success("Changes committed")

	// Create and push tag
	console.Printf("Creating tag...\n")
	if err := gitCreateTag(runner, ctx, newTag); err != nil {
		console.Error(fmt.Sprintf("Error creating tag: %s", err))
		os.Exit(1)
	}
	console.Success(fmt.Sprintf("Tag %s created", newTag))

	// Push commit and tag
	console.Printf("Pushing changes...\n")
	if err := gitPushCommit(runner, ctx); err != nil {
		console.Error(fmt.Sprintf("Error pushing commit: %s", err))
		os.Exit(1)
	}

	if err := gitPushTag(runner, ctx, newTag); err != nil {
		console.Error(fmt.Sprintf("Error pushing tag: %s", err))
		os.Exit(1)
	}

	console.Success(fmt.Sprintf("Tag %s created and pushed successfully!", newTag))
	console.Printf("GoReleaser pipeline should be triggered automatically.\n")
}

func deleteTag(console *printer.Console, reader *bufio.Reader) {
	runner := executil.DefaultRunner()
	ctx := context.Background()

	// Get all tags
	tags, err := getAllTags(runner, ctx)
	if err != nil {
		console.Error(fmt.Sprintf("Error getting tags: %s", err))
		os.Exit(1)
	}

	if len(tags) == 0 {
		console.Warn("No tags found in repository")
		return
	}

	// Display tags
	console.Printf("\nAvailable tags:\n\n")
	for i, tag := range tags {
		console.Printf("  %d) %s\n", i+1, tag)
	}
	console.Printf("\n")

	// Get user input
	console.Printf("Enter tag number to delete (or 'cancel'): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		console.Error(fmt.Sprintf("Error reading input: %s", err))
		os.Exit(1)
	}

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "cancel" || input == "c" {
		console.Printf("Cancelled.\n")
		return
	}

	// Parse selection
	var selectedTag string
	var selection int
	if _, err := fmt.Sscanf(input, "%d", &selection); err != nil {
		console.Error("Invalid input")
		os.Exit(1)
	}

	if selection < 1 || selection > len(tags) {
		console.Error("Invalid selection")
		os.Exit(1)
	}

	selectedTag = tags[selection-1]

	// Confirm deletion
	console.Warn(fmt.Sprintf("Are you sure you want to delete tag '%s'?", selectedTag))
	console.Printf("Confirm (y/n): ")
	confirm, err := reader.ReadString('\n')
	if err != nil {
		console.Error(fmt.Sprintf("Error reading confirmation: %s", err))
		os.Exit(1)
	}

	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" && confirm != "yes" {
		console.Printf("Cancelled.\n")
		return
	}

	// Delete local tag
	if err := gitDeleteLocalTag(runner, ctx, selectedTag); err != nil {
		console.Error(fmt.Sprintf("Error deleting local tag: %s", err))
		os.Exit(1)
	}

	console.Success(fmt.Sprintf("Deleted local tag: %s", selectedTag))

	// Ask if they want to delete remote tag too
	console.Printf("\nDelete remote tag as well? (y/n): ")
	deleteRemote, err := reader.ReadString('\n')
	if err != nil {
		console.Error(fmt.Sprintf("Error reading input: %s", err))
		os.Exit(1)
	}

	deleteRemote = strings.TrimSpace(strings.ToLower(deleteRemote))
	if deleteRemote == "y" || deleteRemote == "yes" {
		if err := gitDeleteRemoteTag(runner, ctx, selectedTag); err != nil {
			console.Error(fmt.Sprintf("Error deleting remote tag: %s", err))
			os.Exit(1)
		}
		console.Success(fmt.Sprintf("Deleted remote tag: %s", selectedTag))
	}

	console.Success("Tag deletion completed successfully!")
}

func getLatestTag(runner executil.Runner, ctx context.Context) string {
	// First, fetch all tags from remote to ensure we have the latest
	fetchCmd := executil.NewCommand("git", "fetch", "--tags")
	_, _ = runner.Run(ctx, fetchCmd) // Ignore errors as repo might not have remote

	// Get the latest tag
	var stdout bytes.Buffer
	describeCmd := executil.NewCommand("git", "describe", "--tags", "--abbrev=0")
	describeCmd.Stdout = &stdout

	_, err := runner.Run(ctx, describeCmd)
	if err != nil {
		// No tags exist yet, start from v0.0.0
		return "v0.0.0"
	}
	return strings.TrimSpace(stdout.String())
}

func getAllTags(runner executil.Runner, ctx context.Context) ([]string, error) {
	var stdout bytes.Buffer
	cmd := executil.NewCommand("git", "tag", "--sort=-version:refname")
	cmd.Stdout = &stdout

	_, err := runner.Run(ctx, cmd)
	if err != nil {
		return nil, err
	}

	tagsStr := strings.TrimSpace(stdout.String())
	if tagsStr == "" {
		return []string{}, nil
	}

	return strings.Split(tagsStr, "\n"), nil
}

func gitCreateTag(runner executil.Runner, ctx context.Context, tag string) error {
	cmd := executil.NewCommand("git", "tag", tag)
	result, err := runner.Run(ctx, cmd)
	if err != nil {
		return fmt.Errorf("%s: %s", err, result.Output)
	}
	return nil
}

func gitPushTag(runner executil.Runner, ctx context.Context, tag string) error {
	cmd := executil.NewCommand("git", "push", "origin", tag)
	result, err := runner.Run(ctx, cmd)
	if err != nil {
		return fmt.Errorf("%s: %s", err, result.Output)
	}
	return nil
}

func gitDeleteLocalTag(runner executil.Runner, ctx context.Context, tag string) error {
	cmd := executil.NewCommand("git", "tag", "-d", tag)
	result, err := runner.Run(ctx, cmd)
	if err != nil {
		return fmt.Errorf("%s: %s", err, result.Output)
	}
	return nil
}

func gitDeleteRemoteTag(runner executil.Runner, ctx context.Context, tag string) error {
	cmd := executil.NewCommand("git", "push", "origin", "--delete", tag)
	result, err := runner.Run(ctx, cmd)
	if err != nil {
		return fmt.Errorf("%s: %s", err, result.Output)
	}
	return nil
}

func generateVersionFile(console *printer.Console, version string) error {
	fs := afero.NewOsFs()
	tracker := manifest.NewTracker()
	gen := scaffold.New(fs, tracker, console)

	// Generate simplified version file content
	content := fmt.Sprintf(`package version

const Version = %q
`, version)

	return gen.Code("internal/version/version.go", content)
}

func runGoGenerate(runner executil.Runner, ctx context.Context) error {
	cmd := executil.NewCommand("go", "generate", "./...")
	result, err := runner.Run(ctx, cmd)
	if err != nil {
		return fmt.Errorf("go generate failed: %s: %s", err, result.Output)
	}
	return nil
}

func gitCommitAllChanges(runner executil.Runner, ctx context.Context, tag string) error {
	// Add all changes (generated files + version file)
	addCmd := executil.NewCommand("git", "add", "-A")
	result, err := runner.Run(ctx, addCmd)
	if err != nil {
		return fmt.Errorf("git add failed: %s: %s", err, result.Output)
	}

	// Check if there are any changes to commit
	statusCmd := executil.NewCommand("git", "diff", "--cached", "--quiet")
	_, err = runner.Run(ctx, statusCmd)
	if err == nil {
		// No changes to commit
		return nil
	}

	// Commit with a message
	commitMsg := fmt.Sprintf("chore: Updating version to %s", tag)
	commitCmd := executil.NewCommand("git", "commit", "-m", commitMsg)
	result, err = runner.Run(ctx, commitCmd)
	if err != nil {
		return fmt.Errorf("git commit failed: %s: %s", err, result.Output)
	}

	return nil
}

func gitPushCommit(runner executil.Runner, ctx context.Context) error {
	cmd := executil.NewCommand("git", "push")
	result, err := runner.Run(ctx, cmd)
	if err != nil {
		return fmt.Errorf("%s: %s", err, result.Output)
	}
	return nil
}
