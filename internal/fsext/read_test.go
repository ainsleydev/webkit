package fsext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromEmbed(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		name     string
		filename string
		want     string
		wantErr  bool
	}{
		"Existing File": {
			filename: "testdata/one.txt",
			want:     "hello world\n",
			wantErr:  false,
		},
		"Non-existent File": {
			filename: "testdata/missing.txt",
			want:     "",
			wantErr:  true,
		},
		"Nested File": {
			filename: "testdata/nested/nested.txt",
			want:     "nested file\n",
			wantErr:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ReadFromEmbed(testFS, test.filename)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMustReadFromEmbed(t *testing.T) {
	t.Parallel()

	t.Run("existing file", func(t *testing.T) {
		t.Parallel()

		got := MustReadFromEmbed(testFS, "testdata/one.txt")
		want := "hello world\n"
		assert.Equal(t, want, got)
	})

	t.Run("Panics on Missing File", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			MustReadFromEmbed(testFS, "testdata/missing.txt")
		})
	})
}
