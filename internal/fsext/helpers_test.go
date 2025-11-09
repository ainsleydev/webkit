package fsext

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		setup func(afero.Fs)
		path  string
		want  bool
	}{
		"File exists": {
			setup: func(fs afero.Fs) {
				_ = afero.WriteFile(fs, "test.txt", []byte("content"), 0o644)
			},
			path: "test.txt",
			want: true,
		},
		"File does not exist": {
			setup: func(fs afero.Fs) {},
			path:  "missing.txt",
			want:  false,
		},
		"Directory exists": {
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll("testdir", 0o755)
			},
			path: "testdir",
			want: true,
		},
		"Directory does not exist": {
			setup: func(fs afero.Fs) {},
			path:  "missing-dir",
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			test.setup(fs)

			got := Exists(fs, test.path)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestDirExists(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		setup func(afero.Fs)
		path  string
		want  bool
	}{
		"Directory exists": {
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll("testdir", 0o755)
			},
			path: "testdir",
			want: true,
		},
		"Directory does not exist": {
			setup: func(fs afero.Fs) {},
			path:  "missing-dir",
			want:  false,
		},
		"File exists but not a directory": {
			setup: func(fs afero.Fs) {
				_ = afero.WriteFile(fs, "test.txt", []byte("content"), 0o644)
			},
			path: "test.txt",
			want: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			test.setup(fs)

			got := DirExists(fs, test.path)
			assert.Equal(t, test.want, got)
		})
	}
}
