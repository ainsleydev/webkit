package sops

// Provider abstracts different key management strategies for
// SOPS encryption. Implementations provide the necessary CLI
// arguments and environment variables for SOPS to execute.
type Provider interface {
	// EncryptArgs returns the CLI args needed for encryption.
	// e.g. ["--age", "age1abc..."] or ["--kms", "arn:aws:kms:..."]
	EncryptArgs() ([]string, error)

	// DecryptArgs returns the CLI args needed for decryption.
	// e.g. ["--age", "age1abc..."] or ["--kms", "arn:aws:kms:..."]
	DecryptArgs() ([]string, error)

	// Environment returns environment variables needed for SOPS operations
	// e.g., ["SOPS_AGE_KEY=AGE-SECRET-KEY-1..."]
	Environment() map[string]string
}
