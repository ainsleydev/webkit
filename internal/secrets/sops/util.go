package sops

import "bytes"

// IsContentEncrypted checks if file contents contain SOPS encryption markers.
func IsContentEncrypted(content []byte) bool {
	// Check for SOPS metadata section
	return bytes.Contains(content, []byte("sops:")) ||
		bytes.Contains(content, []byte("ENC["))
}
