package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ainsleydev/webkit/internal/printer"
)

func main() {
	var (
		version string
		repos   string
	)

	flag.StringVar(&version, "version", "", "WebKit version to update to (e.g., 0.1.1)")
	flag.StringVar(&repos, "repos", "[]", "JSON array of repositories to update")
	flag.Parse()

	p := printer.New(os.Stdout)

	if version == "" {
		p.Error("--version flag is required")
		os.Exit(1)
	}

	// Parse repos JSON
	var repoList []string
	if err := json.Unmarshal([]byte(repos), &repoList); err != nil {
		p.Error(fmt.Sprintf("Failed to parse repos JSON: %v", err))
		os.Exit(1)
	}

	if len(repoList) == 0 {
		p.Info("No repositories to update")
		return
	}

	p.Info(fmt.Sprintf("Updating repositories to WebKit v%s", version))
	p.LineBreak()

	for _, repo := range repoList {
		if err := updateRepo(p, repo, version); err != nil {
			p.Error(fmt.Sprintf("Failed to update %s: %v", repo, err))
			continue
		}
	}
}

func updateRepo(p *printer.Printer, repo, version string) error {
	p.Info(fmt.Sprintf("Processing %s...", repo))

	repoDir := filepath.Base(repo)
	ghToken := os.Getenv("GH_TOKEN")

	// Clone the repository
	if err := runCommand("gh", "repo", "clone", repo, repoDir, "--", "--depth=1"); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	// Change to repo directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	defer os.RemoveAll(repoDir)

	if err := os.Chdir(repoDir); err != nil {
		return fmt.Errorf("chdir failed: %w", err)
	}

	// Configure git remote
	remoteURL := fmt.Sprintf("https://x-access-token:%s@github.com/%s.git", ghToken, repo)
	if err := runCommand("git", "remote", "set-url", "origin", remoteURL); err != nil {
		return fmt.Errorf("git remote set-url failed: %w", err)
	}

	// Check if app.json exists and has webkit_version
	if _, err := os.Stat("app.json"); os.IsNotExist(err) {
		p.Warn("No app.json found, skipping...")
		return nil
	}

	// Read current version
	appJSON, err := os.ReadFile("app.json")
	if err != nil {
		return fmt.Errorf("read app.json failed: %w", err)
	}

	var appData map[string]interface{}
	if err := json.Unmarshal(appJSON, &appData); err != nil {
		return fmt.Errorf("parse app.json failed: %w", err)
	}

	currentVersion, ok := appData["webkit_version"].(string)
	if !ok {
		p.Warn("No webkit_version field in app.json, skipping...")
		return nil
	}

	currentVersion = strings.TrimPrefix(currentVersion, "v")
	if currentVersion == version {
		p.Success(fmt.Sprintf("Already on v%s, skipping...", version))
		return nil
	}

	p.Info(fmt.Sprintf("Updating from v%s to v%s", currentVersion, version))

	branchName := fmt.Sprintf("chore/update-webkit-v%s", version)

	// Close stale PRs
	p.Info("Checking for stale webkit update PRs...")
	if err := closeStalePRs(repo, branchName); err != nil {
		p.Warn(fmt.Sprintf("Could not close stale PRs: %v", err))
	}

	// Check if PR already exists
	output, _ := exec.Command("gh", "pr", "list", "--repo", repo, "--head", branchName, "--json", "number", "--jq", ".[0].number").Output()
	if len(strings.TrimSpace(string(output))) > 0 {
		p.Success(fmt.Sprintf("PR already exists for v%s", version))
		return nil
	}

	// Create new branch
	if err := runCommand("git", "checkout", "-b", branchName); err != nil {
		return fmt.Errorf("git checkout failed: %w", err)
	}

	// Install webkit
	p.Info(fmt.Sprintf("Installing webkit v%s...", version))
	installCmd := fmt.Sprintf("curl -sSL https://raw.githubusercontent.com/ainsleydev/webkit/main/bin/install.sh | VERSION=\"v%s\" sh", version)
	if err := runShell(installCmd); err != nil {
		return fmt.Errorf("webkit install failed: %w", err)
	}

	// Ensure webkit is in PATH
	homeDir, _ := os.UserHomeDir()
	localBin := filepath.Join(homeDir, ".local", "bin")
	os.Setenv("PATH", localBin+":"+os.Getenv("PATH"))

	// Run webkit update
	p.Info("Running webkit update...")
	if err := runCommand("webkit", "update"); err != nil {
		return fmt.Errorf("webkit update failed: %w", err)
	}

	// Check for changes
	if !hasChanges() {
		p.Warn("No changes detected after webkit update, skipping...")
		return nil
	}

	// Commit changes
	// Use environment variables for bot identity, defaulting to ainsleydev-bot
	botName := getEnvOrDefault("GIT_AUTHOR_NAME", "ainsleydev-bot[bot]")
	botEmail := getEnvOrDefault("GIT_AUTHOR_EMAIL", "175332+ainsleydev-bot[bot]@users.noreply.github.com")

	if err := runCommand("git", "config", "user.name", botName); err != nil {
		return fmt.Errorf("git config name failed: %w", err)
	}
	if err := runCommand("git", "config", "user.email", botEmail); err != nil {
		return fmt.Errorf("git config email failed: %w", err)
	}
	if err := runCommand("git", "add", "."); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}
	if err := runCommand("git", "commit", "-m", fmt.Sprintf("chore: Update WebKit to v%s", version)); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	// Delete remote branch if exists
	runCommand("git", "push", "origin", "--delete", branchName) // Ignore errors

	// Push new branch
	if err := runCommand("git", "push", "origin", branchName); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	// Create PR
	prBody := fmt.Sprintf(`This automated PR updates the WebKit version from v%s to v%s.

## Testing
- [ ] Review the generated files for any unexpected changes
- [ ] Check for any breaking changes in the [release notes](https://github.com/ainsleydev/webkit/releases/tag/v%s)
- [ ] Test the application locally after merging

🤖 Generated with ainsley.dev bot`, currentVersion, version, version)

	if err := runCommand("gh", "pr", "create",
		"--repo", repo,
		"--head", branchName,
		"--title", fmt.Sprintf("chore: Update WebKit to v%s", version),
		"--body", prBody,
		"--base", "main"); err != nil {
		return fmt.Errorf("gh pr create failed: %w", err)
	}

	p.Success("Pull request created successfully")
	p.LineBreak()

	return nil
}

func closeStalePRs(repo, currentBranch string) error {
	output, err := exec.Command("gh", "pr", "list",
		"--repo", repo,
		"--search", "chore/update-webkit-v in:title",
		"--state", "open",
		"--json", "number,headRefName").Output()
	if err != nil {
		return err
	}

	var prs []struct {
		Number      int    `json:"number"`
		HeadRefName string `json:"headRefName"`
	}

	if err := json.Unmarshal(output, &prs); err != nil {
		return err
	}

	for _, pr := range prs {
		if pr.HeadRefName != currentBranch && strings.HasPrefix(pr.HeadRefName, "chore/update-webkit-v") {
			runCommand("gh", "pr", "close", fmt.Sprintf("%d", pr.Number),
				"--repo", repo,
				"--comment", "Closing this PR as it's superseded by a newer WebKit version update.")
			runCommand("git", "push", "origin", "--delete", pr.HeadRefName) // Ignore errors
		}
	}

	return nil
}

func hasChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, _ := cmd.Output()
	return len(strings.TrimSpace(string(output))) > 0
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runShell(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
