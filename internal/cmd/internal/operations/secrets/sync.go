package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
)

var SyncCmd = &cli.Command{
	Name:        "sync",
	Usage:       "Sync secret placeholders from app.json",
	Description: "Reads app.json and adds placeholder entries for all secrets with source: 'sops'",
	Action:      cmdtools.Wrap(Sync),
}

// Sync adds missing secret placeholders to SOPS files based on app.json.
// This command reads environment variables with source: "sops" and ensures
// corresponding placeholder entries exist in the appropriate secret files.
// Only works with unencrypted files.
func Sync(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	results, err := secrets.Sync(secrets.SyncConfig{
		FS:     input.FS,
		AppDef: app,
	})
	if err != nil {
		return fmt.Errorf("syncing secrets: %w", err)
	}

	if len(results.Files) == 0 {
		fmt.Println("No secrets with source: 'sops' found in app.json")
		return nil
	}

	fmt.Println("Syncing secrets from app.json...")

	for _, file := range results.Files {
		printFileResult(file)
	}

	// Print summary
	printSummary(results)

	// Return error if any files had errors
	if results.HasErrors() {
		return fmt.Errorf("some files had errors during sync")
	}

	return nil
}

// printFileResult handles output for a single file
func printFileResult(file secrets.SyncResult) {
	fmt.Printf("✓ %s\n", file.FilePath)

	if file.Error != nil {
		fmt.Printf("  ✗ Error: %v\n", file.Error)
		fmt.Println()
		return
	}

	if file.WasMissing {
		fmt.Println("  ⚠ File does not exist")
		fmt.Println()
		return
	}

	if file.WasEncrypted {
		fmt.Println("  ⚠ Skipped (file is encrypted)")
		fmt.Println()
		return
	}

	// Print added secrets
	if len(file.AddedSecrets) > 0 {
		for _, secret := range file.AddedSecrets {
			appList := strings.Join(secret.AppNames, ", ")
			fmt.Printf("  • Added %s (used by: %s)\n", secret.Key, appList)
		}
	}

	// Print skipped count if any
	if file.Skipped > 0 {
		fmt.Printf("  • Skipped %d existing secret%s\n", file.Skipped, pluralize(file.Skipped))
	}

	// If nothing was added or skipped
	if file.Added == 0 && file.Skipped == 0 {
		fmt.Println("  • No secrets to sync")
	}

	fmt.Println()
}

// printSummary prints overall sync summary
func printSummary(results *secrets.SyncResults) {
	separator := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	fmt.Println(separator)
	fmt.Printf("Summary: %d secret%s added, %d skipped",
		results.TotalAdded(), pluralize(results.TotalAdded()),
		results.TotalSkipped())

	if results.MissingCount() > 0 {
		fmt.Printf(", %d file%s missing", results.MissingCount(), pluralize(results.MissingCount()))
	}
	if results.EncryptedCount() > 0 {
		fmt.Printf(", %d encrypted file%s skipped", results.EncryptedCount(), pluralize(results.EncryptedCount()))
	}

	fmt.Printf("\n%s\n", separator)

	// Helpful hints
	if results.MissingCount() > 0 {
		fmt.Println("\n💡 Run 'webkit scaffold secrets' to create missing files")
	}
	if results.TotalAdded() > 0 {
		fmt.Println("\n💡 Replace placeholder values (REPLACE_ME_*) with actual secrets")
		fmt.Println("💡 Run 'webkit secrets encrypt' to encrypt files before committing")
	}
}

// pluralize returns "s" if count != 1, empty string otherwise
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
