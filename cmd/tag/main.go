package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/ainsleydev/webkit/internal/printer"
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
	// Get latest tag
	latestTag, err := getLatestTag()
	if err != nil {
		console.Error(fmt.Sprintf("Error getting latest tag: %s", err))
		os.Exit(1)
	}

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

	// Create and push tag
	if err := gitCreateTag(newTag); err != nil {
		console.Error(fmt.Sprintf("Error creating tag: %s", err))
		os.Exit(1)
	}

	if err := gitPushTag(newTag); err != nil {
		console.Error(fmt.Sprintf("Error pushing tag: %s", err))
		os.Exit(1)
	}

	console.Success(fmt.Sprintf("Tag %s created and pushed successfully!", newTag))
	console.Printf("GoReleaser pipeline should be triggered automatically.\n")
}

func deleteTag(console *printer.Console, reader *bufio.Reader) {
	// Get all tags
	tags, err := getAllTags()
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
	if err := gitDeleteLocalTag(selectedTag); err != nil {
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
		if err := gitDeleteRemoteTag(selectedTag); err != nil {
			console.Error(fmt.Sprintf("Error deleting remote tag: %s", err))
			os.Exit(1)
		}
		console.Success(fmt.Sprintf("Deleted remote tag: %s", selectedTag))
	}

	console.Success("Tag deletion completed successfully!")
}

func getLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		// No tags exist yet, start from v0.0.0
		return "v0.0.0", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func getAllTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "--sort=-version:refname")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	tagsStr := strings.TrimSpace(string(output))
	if tagsStr == "" {
		return []string{}, nil
	}

	return strings.Split(tagsStr, "\n"), nil
}

func gitCreateTag(tag string) error {
	cmd := exec.Command("git", "tag", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func gitPushTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func gitDeleteLocalTag(tag string) error {
	cmd := exec.Command("git", "tag", "-d", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func gitDeleteRemoteTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", "--delete", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}
