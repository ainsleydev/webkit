package cmd

import (
	"bytes"
	"context"
	"fmt"
	goFs "io/fs"
	"path/filepath"

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
	fs := input.FS

	before := map[string][]byte{}
	err := afero.Walk(fs, ".", func(path string, info goFs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		b, _ := afero.ReadFile(fs, path)
		before[path] = b
		return nil
	})
	if err != nil {
		return fmt.Errorf("snapshot failed: %w", err)
	}

	// run update
	if err := update(ctx, input); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	var changed []string
	_ = afero.Walk(fs, ".", func(path string, info goFs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		after, _ := afero.ReadFile(fs, path)
		beforeBytes, ok := before[path]
		if !ok {
			changed = append(changed, path+" (new)")
			return nil
		}
		if !bytes.Equal(beforeBytes, after) {
			changed = append(changed, filepath.Clean(path))
		}
		return nil
	})

	if len(changed) > 0 {
		fmt.Println("⚠️  Drift detected! The following files differ:")
		for _, f := range changed {
			fmt.Println(" -", f)
		}
		fmt.Println("\nRun `webkit update` and commit the changes if correct.")
		return fmt.Errorf("drift detected")
	}

	fmt.Println("✅ No drift detected — files are up to date.")
	return nil
}
