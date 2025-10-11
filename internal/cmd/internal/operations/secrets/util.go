package secrets

import (
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
)

func getSopsClient() (sops.EncrypterDecrypter, error) {
	prov, err := age.NewProvider()
	if err != nil {
		return nil, err
	}
	return sops.NewClient(prov), nil
}
