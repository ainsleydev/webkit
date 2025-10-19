package manifest

import "github.com/spf13/afero"

// DriftReason represents why a file has drifted
type DriftReason int

const (
	// DriftReasonModified indicates the file content has changed.
	DriftReasonModified DriftReason = iota
	// DriftReasonDeleted indicates the file no longer exists.
	DriftReasonDeleted
)

var driftReasonStrings = map[DriftReason]string{
	DriftReasonModified: "modified",
	DriftReasonDeleted:  "deleted",
}

// String implements fmt.Stringer on the DriftReason.
func (r DriftReason) String() string {
	if s, ok := driftReasonStrings[r]; ok {
		return s
	}
	return "unknown"
}

// DriftedFile represents a file that has drifted from its tracked state.
type DriftedFile struct {
	Path   string
	Reason DriftReason
}

// DetectDrift checks if files have been manually modified and returns a list
// of files that have changed.
//
// If the list is empty, no changes have been made.
func DetectDrift(fs afero.Fs, manifest *Manifest) []DriftedFile {
	var drifted []DriftedFile

	for path, entry := range manifest.Files {
		// Don't try and hash stuff that's managed by the
		// user and not WebKit.
		if entry.ScaffoldMode {
			continue
		}

		// Check if the file exists.
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			// File might be deleted or moved
			drifted = append(drifted, DriftedFile{
				Path:   path,
				Reason: DriftReasonDeleted,
			})
			continue
		}

		// Check if content has changed
		currentHash := HashContent(data)
		if currentHash != entry.Hash {
			drifted = append(drifted, DriftedFile{
				Path:   path,
				Reason: DriftReasonModified,
			})
		}
	}

	return drifted
}
