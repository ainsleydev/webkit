package docs

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

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

			got, err := LoadGenFile(fs, test.file)

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

func TestHasAppType(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def     *appdef.Definition
		appType appdef.AppType
		want    bool
	}{
		"Has Payload app": {
			def: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "cms", Type: appdef.AppTypePayload},
				},
			},
			appType: appdef.AppTypePayload,
			want:    true,
		},
		"Has SvelteKit app": {
			def: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "web", Type: appdef.AppTypeSvelteKit},
				},
			},
			appType: appdef.AppTypeSvelteKit,
			want:    true,
		},
		"Does not have Payload app": {
			def: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "api", Type: appdef.AppTypeGoLang},
				},
			},
			appType: appdef.AppTypePayload,
			want:    false,
		},
		"Nil definition": {
			def:     nil,
			appType: appdef.AppTypePayload,
			want:    false,
		},
		"Empty apps": {
			def:     &appdef.Definition{Apps: []appdef.App{}},
			appType: appdef.AppTypePayload,
			want:    false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := HasAppType(test.def, test.appType)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetAppsByType(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		def     *appdef.Definition
		appType appdef.AppType
		want    int
	}{
		"Multiple Payload apps": {
			def: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "cms", Type: appdef.AppTypePayload},
					{Name: "admin", Type: appdef.AppTypePayload},
					{Name: "web", Type: appdef.AppTypeSvelteKit},
				},
			},
			appType: appdef.AppTypePayload,
			want:    2,
		},
		"Single SvelteKit app": {
			def: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "web", Type: appdef.AppTypeSvelteKit},
					{Name: "api", Type: appdef.AppTypeGoLang},
				},
			},
			appType: appdef.AppTypeSvelteKit,
			want:    1,
		},
		"No matching apps": {
			def: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "api", Type: appdef.AppTypeGoLang},
				},
			},
			appType: appdef.AppTypePayload,
			want:    0,
		},
		"Nil definition": {
			def:     nil,
			appType: appdef.AppTypePayload,
			want:    0,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := GetAppsByType(test.def, test.appType)
			assert.Len(t, got, test.want)
		})
	}
}
