package sops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsContentEncrypted(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input []byte
		want  bool
	}{
		"Empty Content": {
			input: []byte(""),
			want:  false,
		},
		"No SOPS Markers": {
			input: []byte("plain text without encryption markers"),
			want:  false,
		},
		"Contains SOPS Metadata": {
			input: []byte("some content\nsops:\n  kms: ..."),
			want:  true,
		},
		"Contains ENC Marker": {
			input: []byte("some content ENC[AES256_GCM,data...] more text"),
			want:  true,
		},
		"Both SOPS and ENC": {
			input: []byte("sops:\nENC[AES256_GCM,data...]"),
			want:  true,
		},
		"SOPS Lowercase Only": {
			input: []byte("sops: something else"),
			want:  true,
		},
		"ENC Lowercase Only": {
			input: []byte("enc[data]"),
			want:  false, // Function is case-sensitive
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
