package docsutil

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplate_String(t *testing.T) {
	t.Parallel()

	got := CodeStyleTemplate.String()
	assert.Equal(t, "CODE_STYLE.md", got)
	assert.IsType(t, "", got)
}

func TestLoadGenFile(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		setup   func(fs afero.Fs)
		file    string
		want    string
		wantErr bool
	}{
		"File exists": {
			setup: func(fs afero.Fs) {
				err := fs.MkdirAll(genDocsDir, 0755)
				require.NoError(t, err)
				err = afero.WriteFile(fs, filepath.Join(genDocsDir, "CODE_STYLE.md"), []byte("# Code Style"), 0644)
				require.NoError(t, err)
			},
			file:    "CODE_STYLE.md",
			want:    "# Code Style",
			wantErr: false,
		},
		"File does not exist": {
			setup:   func(fs afero.Fs) {},
			file:    "MISSING.md",
			want:    "",
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			test.setup(fs)

			got, err := LoadGenFile(fs, Template(test.file))

			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestLoadCustomContent(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		setup   func(fs afero.Fs)
		want    string
		wantErr bool
	}{
		"Custom content exists": {
			setup: func(fs afero.Fs) {
				err := fs.MkdirAll(customDocsDir, 0755)
				require.NoError(t, err)
				err = afero.WriteFile(fs, filepath.Join(customDocsDir, agentsFilename), []byte("## WebKit\n\nCustom content."), 0644)
				require.NoError(t, err)
			},
			want:    "## WebKit\n\nCustom content.",
			wantErr: false,
		},
		"Custom content does not exist": {
			setup:   func(fs afero.Fs) {},
			want:    "",
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			test.setup(fs)

			got, err := LoadCustomContent(fs)

			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}
