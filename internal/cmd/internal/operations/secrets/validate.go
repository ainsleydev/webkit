package secrets

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
)

var ValidateCmd = &cli.Command{
	Name:        "validate",
	Usage:       "Validate that all secrets from app.json exist in secret files",
	Description: "Ensures every secret referenced in app.json has a corresponding entry in SOPS files",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "check-orphans",
			Usage:   "Report keys in SOPS files not referenced in app.json",
			Aliases: []string{"o"},
		},
		&cli.BoolFlag{
			Name:    "allow-encrypted",
			Usage:   "Attempt to validate encrypted files (requires SOPS/age access)",
			Aliases: []string{"e"},
		},
	},
	Action: cmdtools.Wrap(Validate),
}

// Validate validates that all secrets referenced in app.json exist in SOPS files
func Validate(ctx context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	checkOrphans := input.Command.Bool("check-orphans")
	allowEncrypted := input.Command.Bool("allow-encrypted")

	fmt.Println("Validating secrets from app.json...")

	result, err := secrets.Validate(secrets.ValidateConfig{
		FS:             input.FS,
		AppDef:         app,
		CheckOrphans:   checkOrphans,
		AllowEncrypted: allowEncrypted,
	})
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	printValidationResults(result, checkOrphans)

	if !result.Valid {
		return fmt.Errorf("validation failed - missing secrets detected")
	}

	return nil
}

// printValidationResults outputs validation results in a user-friendly format
func printValidationResults(result *secrets.ValidationResult, checkOrphans bool) {
	separator := "═══════════════════════════════════════════════════════════"

	fmt.Println()

	// Print file-by-file results
	for _, file := range result.Files {
		printFileValidation(file)
	}

	fmt.Println(separator)

	// Print summary
	if result.Valid {
		fmt.Println("✅ All secrets validated successfully!")

		if len(result.Files) > 0 {
			totalKeys := 0
			for _, file := range result.Files {
				totalKeys += len(file.MissingKeys)
			}
			fmt.Printf("   Validated %d secret file(s)\n", len(result.Files))
		}
	} else {
		fmt.Println("❌ Validation failed!")
		fmt.Printf("   %d missing secret(s) detected\n", len(result.MissingSecrets))

		// Group by environment
		byEnv := make(map[string][]secrets.MissingSecret)
		for _, missing := range result.MissingSecrets {
			byEnv[missing.Environment] = append(byEnv[missing.Environment], missing)
		}

		fmt.Println("\nMissing secrets:")
		for env, secrets := range byEnv {
			fmt.Printf("\n  %s:\n", env)
			for _, secret := range secrets {
				fmt.Printf("    • %s (used by: %s)\n", secret.Key, secret.AppName)
				fmt.Printf("      Expected in: %s\n", secret.ExpectedIn)
			}
		}
	}

	// Print orphaned keys if checked
	if checkOrphans && len(result.OrphanedKeys) > 0 {
		fmt.Println("\n⚠️  Orphaned keys (in SOPS files but not in app.json):")
		byFile := make(map[string][]string)
		for _, orphan := range result.OrphanedKeys {
			byFile[orphan.FilePath] = append(byFile[orphan.FilePath], orphan.Key)
		}
		for file, keys := range byFile {
			fmt.Printf("\n  %s:\n", file)
			for _, key := range keys {
				fmt.Printf("    • %s\n", key)
			}
		}
	}

	fmt.Println(separator)

	// Helpful hints
	if !result.Valid {
		fmt.Println("\n💡 Run 'webkit secrets sync' to add missing secrets")
		fmt.Println("💡 Run 'webkit secrets validate --check-orphans' to find unused keys")
	}
}

// printFileValidation outputs validation results for a single file
func printFileValidation(file secrets.FileValidation) {
	if !file.Exists {
		fmt.Printf("⚠️  %s - file does not exist\n", file.FilePath)
		return
	}

	if file.Error != nil {
		fmt.Printf("❌ %s - error: %v\n", file.FilePath, file.Error)
		return
	}

	if file.IsEncrypted {
		fmt.Printf("🔒 %s - encrypted (skipped)\n", file.FilePath)
		return
	}

	if len(file.MissingKeys) == 0 {
		fmt.Printf("✅ %s\n", file.FilePath)
	} else {
		fmt.Printf("❌ %s - %d missing key(s)\n", file.FilePath, len(file.MissingKeys))
		for _, key := range file.MissingKeys {
			fmt.Printf("     • %s\n", key)
		}
	}
}
