package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
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
	printer := input.Printer()

	files := syncSecrets(input.FS, app)
	if len(files) == 0 {
		printer.Warn("No secrets with source: 'sops' found in app.json")
		return nil
	}

	added, skipped, missing, encrypted, failed := 0, 0, 0, 0, 0
	for _, f := range files {
		printer.Printf("✓ %s", f.Path)
		switch {
		case f.Error != nil:
			printer.Printf("  ✗ %v\n", f.Error)
			failed++
		case f.Missing:
			printer.Println("  ⚠ Missing file")
			missing++
		case f.Encrypted:
			printer.Println("  ⚠ Encrypted (skipped)")
			encrypted++
		default:
			printer.Printf("  • %d added, %d skipped\n", f.Added, f.Skipped)
			added += f.Added
			skipped += f.Skipped
		}
	}

	printer.Printf("Summary: %d added, %d skipped, %d missing, %d encrypted, %d failed\n",
		added, skipped, missing, encrypted, failed)

	if failed > 0 {
		return fmt.Errorf("sync completed with errors")
	}

	return nil
}

type syncFile struct {
	Path      string
	Added     int
	Skipped   int
	Missing   bool
	Encrypted bool
	Error     error
}

func syncSecrets(fs afero.Fs, def *appdef.Definition) []syncFile {
	type ref struct {
		Key string
		Env env.Environment
	}

	var refs []ref
	def.MergeAllEnvironments().Walk(func(env env.Environment, name string, value appdef.EnvValue) {
		if value.Source == appdef.EnvSourceSOPS {
			refs = append(refs, ref{Key: name, Env: env})
		}
	})

	if len(refs) == 0 {
		return nil
	}

	// Group by environment (file)
	group := map[env.Environment][]string{}
	for _, r := range refs {
		group[r.Env] = append(group[r.Env], r.Key)
	}

	var results []syncFile
	for env, keys := range group {
		path := secrets.FilePathFromEnv(env)
		results = append(results, processSyncFile(fs, path, keys))
	}

	return results
}

// processSyncFile processes a single secret file by adding missing placeholders.
// It checks if the file exists, is encrypted, and adds any missing secret keys.
func processSyncFile(fs afero.Fs, path string, keys []string) syncFile {
	result := syncFile{Path: path}

	content, err := afero.ReadFile(fs, path)
	if err != nil {
		result.Missing = true
		return result
	}

	if sops.IsContentEncrypted(content) {
		result.Encrypted = true
		return result
	}

	data := map[string]any{}
	if err = yaml.Unmarshal(content, &data); err != nil {
		result.Error = fmt.Errorf("invalid YAML: %w", err)
		return result
	}

	// Process each secret reference and check if they exist
	// in the file; if they don't add a placeholder.
	var sb strings.Builder
	for _, key := range keys {
		if _, ok := data[key]; ok {
			result.Skipped++
			continue
		}
		sb.WriteString(fmt.Sprintf("%s: \"REPLACE_ME_%s\"\n", key, strings.ToUpper(key)))
		result.Added++
	}

	// Write back to the file if any of the secrets
	// need scaffolding to the file.
	if result.Added > 0 {
		content = append(content, []byte(sb.String())...)
		err = afero.WriteFile(fs, path, content, 0644)
		if err != nil {
			result.Error = err
		}
	}

	return result
}
