package manifest

import (
	"fmt"

	"github.com/spf13/afero"
)

// DriftReason represents the kind of drift detected
type DriftReason int

// DriftReason constants.
const (
	DriftReasonModified DriftReason = iota // File was manually edited
	DriftReasonDeleted                     // File should be removed
	DriftReasonOutdated                    // app.json changed, needs regen
	DriftReasonNew                         // File should exist but doesn't
)

// driftReasonStrings maps DriftReason values to their string representations.
var driftReasonStrings = map[DriftReason]string{
	DriftReasonModified: "modified",
	DriftReasonDeleted:  "deleted",
	DriftReasonOutdated: "outdated",
	DriftReasonNew:      "new",
}

// String implements fmt.Stringer on the DriftReason.
func (d DriftReason) String() string {
	if s, ok := driftReasonStrings[d]; ok {
		return s
	}
	return fmt.Sprintf("unknown(%d)", int(d))
}

// FilterEntries plucks entries by the selected DriftReason.
func (d DriftReason) FilterEntries(filtered []DriftEntry) []DriftEntry {
	var out []DriftEntry
	for _, entry := range filtered {
		if entry.Type == d {
			out = append(out, entry)
		}
	}
	return out
}

// DriftEntry contains information about a single drifted file
type DriftEntry struct {
	Path      string      // File path
	Type      DriftReason // Type of drift
	Source    string      // What in app.json caused this
	Generator string      // Which generator created it
}

// DetectDrift compares actual files on disk against what should be generated.
// It detects both manual modifications and source drift (app.json changes).
//
// actualFS: The real filesystem
// expectedFS: In-memory filesystem with freshly generated files from current app.json
//
// Returns all drift: manual edits, outdated files, missing files, and orphaned files.
func DetectDrift(actualFS, expectedFS afero.Fs) ([]DriftEntry, error) {
	var drifted []DriftEntry

	// Load what was previously generated
	actualManifest, err := Load(actualFS)
	if err != nil {
		// No manifest means nothing to compare
		return nil, err
	}

	// Load what should be generated now
	expectedManifest, err := Load(expectedFS)
	if err != nil {
		return nil, err
	}

	// Track which files we've checked
	checkedFiles := make(map[string]bool)

	// Check all files that should exist
	for path, expectedEntry := range expectedManifest.Files {
		checkedFiles[path] = true

		// Skip user-managed files
		if expectedEntry.ScaffoldMode {
			continue
		}

		// Read expected content
		expectedData, err := afero.ReadFile(expectedFS, path)
		if err != nil {
			continue
		}

		// Read actual content
		actualData, err := afero.ReadFile(actualFS, path)
		if err != nil {
			// File should exist but doesn't
			drifted = append(drifted, DriftEntry{
				Path:      path,
				Type:      DriftReasonNew,
				Source:    expectedEntry.Source,
				Generator: expectedEntry.Generator,
			})
			continue
		}

		// Compare content
		expectedHash := HashContent(expectedData)
		actualHash := HashContent(actualData)

		if expectedHash != actualHash {
			// Determine if this is a manual edit or outdated from app.json change
			driftType := DriftReasonOutdated

			// If file exists in old manifest with same hash, it was manually modified
			if oldEntry, exists := actualManifest.Files[path]; exists {
				if oldEntry.Hash == HashContent(actualData) {
					// File matches what we last generated, so app.json must have changed
					driftType = DriftReasonOutdated
				} else {
					// File doesn't match what we last generated, so user edited it
					driftType = DriftReasonModified
				}
			}

			drifted = append(drifted, DriftEntry{
				Path:      path,
				Type:      driftType,
				Source:    expectedEntry.Source,
				Generator: expectedEntry.Generator,
			})
		}
	}

	// Check for orphaned files (in old manifest but not expected anymore)
	for path, entry := range actualManifest.Files {
		if entry.ScaffoldMode || checkedFiles[path] {
			continue
		}

		// File was generated before but shouldn't be generated now
		_, existsInExpected := expectedManifest.Files[path]
		if !existsInExpected {
			exists, _ := afero.Exists(actualFS, path)
			if exists {
				drifted = append(drifted, DriftEntry{
					Path:      path,
					Type:      DriftReasonDeleted,
					Source:    entry.Source,
					Generator: entry.Generator,
				})
			}
		}
	}

	return drifted, nil
}
