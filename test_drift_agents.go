// +build ignore

// Test script to reproduce AGENTS.md drift issue
// Run with: go run test_drift_agents.go

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/docs"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

func main() {
	fmt.Println("=== Testing AGENTS.md Drift ===\n")

	// Create test app definition
	appDef := &appdef.Definition{
		Project: appdef.Project{
			Name: "test-project",
		},
	}

	// Step 1: Generate AGENTS.md first time
	fmt.Println("Step 1: Generating AGENTS.md (first time)...")
	fs1 := afero.NewMemMapFs()
	tracker1 := manifest.NewTracker()
	input1 := cmdtools.CommandInput{
		FS:          fs1,
		AppDefCache: appDef,
		Manifest:    tracker1,
		BaseDir:     "./",
	}

	err := docs.Agents(context.Background(), input1)
	if err != nil {
		fmt.Printf("Error generating AGENTS.md: %v\n", err)
		os.Exit(1)
	}

	// Save manifest
	err = tracker1.Save(fs1)
	if err != nil {
		fmt.Printf("Error saving manifest: %v\n", err)
		os.Exit(1)
	}

	content1, _ := afero.ReadFile(fs1, "AGENTS.md")
	hash1 := manifest.HashContent(content1)
	fmt.Printf("Generated AGENTS.md with hash: %s\n", hash1)

	// Step 2: Load manifest and check stored hash
	savedManifest, err := manifest.Load(fs1)
	if err != nil {
		fmt.Printf("Error loading manifest: %v\n", err)
		os.Exit(1)
	}
	storedHash := savedManifest.Files["AGENTS.md"].Hash
	fmt.Printf("Stored hash in manifest: %s\n", storedHash)

	if hash1 != storedHash {
		fmt.Printf("❌ BUG: File hash doesn't match stored hash!\n")
		os.Exit(1)
	}
	fmt.Printf("✓ File hash matches stored hash\n\n")

	// Step 3: Generate AGENTS.md second time (simulate drift check)
	fmt.Println("Step 2: Generating AGENTS.md again (simulate drift check)...")
	fs2 := afero.NewMemMapFs()
	tracker2 := manifest.NewTracker()
	input2 := cmdtools.CommandInput{
		FS:          fs2,
		AppDefCache: appDef,
		Manifest:    tracker2,
		BaseDir:     "./",
	}

	err = docs.Agents(context.Background(), input2)
	if err != nil {
		fmt.Printf("Error generating AGENTS.md: %v\n", err)
		os.Exit(1)
	}

	content2, _ := afero.ReadFile(fs2, "AGENTS.md")
	hash2 := manifest.HashContent(content2)
	fmt.Printf("Regenerated AGENTS.md with hash: %s\n", hash2)

	// Step 4: Compare hashes
	if hash1 != hash2 {
		fmt.Printf("\n❌ BUG FOUND: AGENTS.md generation is NON-DETERMINISTIC!\n")
		fmt.Printf("First generation hash:  %s\n", hash1)
		fmt.Printf("Second generation hash: %s\n", hash2)
		fmt.Printf("\nThis will cause false drift warnings!\n")

		// Show diff
		fmt.Println("\n=== Content Diff ===")
		fmt.Printf("Length: %d vs %d\n", len(content1), len(content2))

		// Find first difference
		minLen := len(content1)
		if len(content2) < minLen {
			minLen = len(content2)
		}
		for i := 0; i < minLen; i++ {
			if content1[i] != content2[i] {
				fmt.Printf("First difference at byte %d: %q vs %q\n", i, content1[i], content2[i])
				// Show context
				start := i - 20
				if start < 0 {
					start = 0
				}
				end := i + 20
				if end > minLen {
					end = minLen
				}
				fmt.Printf("Context1: %q\n", content1[start:end])
				fmt.Printf("Context2: %q\n", content2[start:end])
				break
			}
		}

		os.Exit(1)
	}

	fmt.Printf("✓ Both generations produced identical content\n\n")

	// Step 5: Test drift detection
	fmt.Println("Step 3: Testing drift detection...")

	// Copy first filesystem as "actual"
	actualFS := afero.NewMemMapFs()
	afero.WriteFile(actualFS, "AGENTS.md", content1, 0644)
	tracker1.Save(actualFS)

	// Use second filesystem as "expected"
	tracker2.Save(fs2)

	drift, err := manifest.DetectDrift(actualFS, fs2)
	if err != nil {
		fmt.Printf("Error detecting drift: %v\n", err)
		os.Exit(1)
	}

	if len(drift) > 0 {
		fmt.Printf("❌ BUG: Drift detected even though content is identical!\n")
		for _, d := range drift {
			fmt.Printf("  - %s: %s\n", d.Path, d.Type.String())
		}
		os.Exit(1)
	}

	fmt.Printf("✓ No drift detected\n\n")
	fmt.Println("=== All tests passed! ===")
}
