package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var driftCmd = &cli.Command{
	Name:   "drift",
	Usage:  "Detect file drift caused by outdated WebKit templates",
	Action: cmdtools.WrapCommand(driftDetection),
}

func driftDetection(ctx context.Context, input cmdtools.CommandInput) error {
	// Capture modification times before update
	before := map[string]os.FileInfo{}
	err := afero.Walk(input.FS, ".", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			before[path] = info
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to snapshot files: %w", err)
	}

	// Run update (idempotent regeneration)
	if err := update(ctx, input); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// Compare after update
	var changed []string
	_ = afero.Walk(input.FS, ".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if prev, ok := before[path]; !ok {
			changed = append(changed, path+" (new)")
		} else if info.ModTime() != prev.ModTime() || info.Size() != prev.Size() {
			changed = append(changed, path)
		}
		return nil
	})

	if len(changed) > 0 {
		fmt.Println("⚠️  Drift detected! The following files were modified:")
		for _, f := range changed {
			fmt.Println(" -", f)
		}
		fmt.Println("\nRun `webkit update` and commit the changes if correct.")
		return fmt.Errorf("drift detected")
	}

	fmt.Println("✅ No drift detected — files are up to date.")
	return nil
}
