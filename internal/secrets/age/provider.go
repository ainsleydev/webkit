package age

import "github.com/ainsleydev/webkit/internal/secrets/sops"

// Provider implements the SOPS Provider interface for
// age encryption, keys are lazy loaded.
type Provider struct {
	privateKey string
	publicKey  string
}

// NewProvider creates a new age provider by reading keys.
//
// Returns an error if it couldn't extract/read public
// and private keys.
func NewProvider() (*Provider, error) {
	identity, err := ReadIdentity()
	if err != nil {
		return nil, err
	}
	return &Provider{
		privateKey: identity.String(),
		publicKey:  identity.Recipient().String(),
	}, nil
}

var _ sops.Provider = (*Provider)(nil)

func (p *Provider) EncryptArgs() ([]string, error) {
	return []string{"--age", p.publicKey}, nil
}

func (p *Provider) DecryptArgs() ([]string, error) {
	return []string{}, nil // Decrypt uses environment variable
}

func (p *Provider) Environment() map[string]string {
	return map[string]string{
		"SOPS_AGE_KEY": p.privateKey,
	}
}
