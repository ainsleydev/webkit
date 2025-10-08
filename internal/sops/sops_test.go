package sops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEncrypted(t *testing.T) {
	t.Parallel()
	t.Skip()

	encryptedYAML := []byte(`{
	"data": "ENC[AES256_GCM,data:test,iv:test,tag:test,type:str]",
	"sops": {
		"kms": [{"arn": "fake-arn"}],
		"version": "3.7.1"
	}
}`)

	unencryptedYAML := []byte(`some: data
not_sops: nope`)

	emptyFile := []byte(``)

	tt := map[string]struct {
		input []byte
		want  bool
	}{
		"Encrypted File": {
			input: encryptedYAML,
			want:  true,
		},
		"Unencrypted File": {
			input: unencryptedYAML,
			want:  false,
		},
		"Empty File": {
			input: emptyFile,
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := IsContentEncrypted(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
