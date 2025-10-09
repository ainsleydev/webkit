package sops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type (
	// Encrypter encrypts a SOPS file using the CLI with the specified
	// provider. Uses the SOPS CLI, assumes it's installed.
	//
	// Example:
	//
	//	err := sops.Encrypt("secrets/production.yaml")
	Encrypter interface {
		Encrypt(filePath string) error
	}
	// Decrypter decrypts a SOPS file using the CLI with the specified
	// provider. Uses the SOPS cli, assumes it's installed.
	//
	// Example:
	//
	//	err := sops.Decrypter("secrets/production.yaml")
	Decrypter interface {
		Decrypt(filePath string) error
	}
	// EncrypterDecrypter combines encryption and decryption operations.
	EncrypterDecrypter interface {
		Encrypter
		Decrypter
	}
)

// Client executes SOPS operations using a configured provider.
type Client struct {
	provider Provider
	exec     func(ctx context.Context, name string, arg ...string) *exec.Cmd
}

// NewClient creates a SOPS client with the specified provider
func NewClient(provider Provider) *Client {
	return &Client{
		provider: provider,
		exec:     exec.CommandContext,
	}
}

var (
	// ErrAlreadyEncrypted is returned when Encrypt is called on a file
	// that is already encrypted.
	ErrAlreadyEncrypted = errors.New("file is already encrypted")

	// ErrNotEncrypted is returned when Decrypt is called on a file that
	// is not encrypted.
	ErrNotEncrypted = errors.New("file is not encrypted")
)

func (c Client) Decrypt(filePath string) error {
	outStr, err := c.runSopsCommand("--decrypt", "--in-place", filePath)

	if strings.Contains(outStr, "sops metadata not found") {
		return ErrNotEncrypted
	}

	if err != nil {
		return fmt.Errorf("sops decrypt failed: %s: %w", outStr, err)
	}

	return nil
}

func (c Client) Encrypt(filePath string) error {
	encryptArgs, err := c.provider.EncryptArgs()
	if err != nil {
		return err
	}

	args := append([]string{"--encrypt"}, encryptArgs...)
	args = append(args, "--in-place", filePath)

	outStr, err := c.runSopsCommand(args...)
	outStr = strings.TrimSpace(outStr)

	if strings.Contains(outStr, "contains a top-level entry called 'sops'") {
		return ErrAlreadyEncrypted
	}

	if err != nil {
		return fmt.Errorf("sops encrypt failed: %s: %w", outStr, err)
	}

	return nil
}

func (c Client) runSopsCommand(args ...string) (string, error) {
	cmd := c.exec(context.Background(), "sops", args...)
	cmd.Env = os.Environ()

	for k, v := range c.provider.Environment() {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}
