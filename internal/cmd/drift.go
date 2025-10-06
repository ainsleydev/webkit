package cmd

import (
	"bytes"
	"context"
	"fmt"
	goFs "io/fs"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var driftCmd = &cli.Command{
	Name:        "drift",
	Description: "Detect file drift caused by outdated WebKit templates",
	Action:      cmdtools.WrapCommand(driftDetection),
}

func driftDetection(ctx context.Context, input cmdtools.CommandInput) error {
	fs := input.FS

	// Capture all of the files and their contents so we have
	// something to compare it too.
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

	// Run update so we can see if the user has made any changes
	// to the root templates.
	if err = update(ctx, input); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// Check if anything has changed after we've updated.
	var driftFound bool
	_ = afero.Walk(fs, ".", func(path string, info goFs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		after, _ := afero.ReadFile(fs, path)
		beforeBytes, ok := before[path]
		if !ok {
			fmt.Printf("🆕 New file: %s\n", path)
			driftFound = true
			return nil
		}

		if !bytes.Equal(beforeBytes, after) {
			driftFound = true
			//printDiff(path, beforeBytes, after) // 👈 replaces cmp.Diff
			fmt.Println("-------")
			printUnifiedDiff(path, beforeBytes, after)
			//fmt.Printf("\n🔍 Diff for %s:\n%s\n", filepath.Clean(path), cmp.Diff(string(beforeBytes), string(after)))
		}
		return nil
	})

	if driftFound {
		fmt.Println("⚠️  Drift detected! Run `webkit update` and commit changes if correct.")
		return fmt.Errorf("drift detected")
	}

	fmt.Println("✅ No drift detected — files are up to date.")
	return nil
}

//func printDiff(path string, before, after []byte) {
//	dmp := diffmatchpatch.New()
//	diffs := dmp.DiffMain(string(before), string(after), false)
//	fmt.Printf("\n🔍 %s\n%s\n", filepath.Clean(path), dmp.(diffs))
//}

func printUnifiedDiff(path string, before, after []byte) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(before)),
		B:        difflib.SplitLines(string(after)),
		FromFile: path + " (before)",
		ToFile:   path + " (after)",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	fmt.Println(text)
}
