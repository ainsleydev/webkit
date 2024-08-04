package static

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestImage_ImageMarkup(t *testing.T) {
	input := NewImage("/assets/images/cat.jpg")
	got := input.ImageMarkup()
	want := markup.ImageProps{
		URL: "/assets/images/cat.jpg",
	}
	assert.Equal(t, want, got)
}

func TestImage_PictureMarkup(t *testing.T) {
	orig := getDistPath
	defer func() {
		getDistPath = orig
	}()
	getDistPath = func(path string) string {
		wd, err := os.Getwd()
		require.NoError(t, err)
		testdata := strings.Replace(path, "assets", "testdata", 1)
		return filepath.Join(wd, testdata)
	}

	t.Run("OK", func(t *testing.T) {
		input := Image("/assets/images/gopher.jpg")
		got := input.PictureMarkup()

		assert.Equal(t, got.URL, input.String())
		require.Len(t, got.Sources, 14)

		assert.Equal(t, got.Sources[0], markup.ImageProps{URL: "/assets/images/gopher-thumbnail-400x300.avif", IsSource: true, MimeType: markup.ImageMimeTypeAVIF, Width: ptr.IntPtr(400), Height: ptr.IntPtr(300)})
		assert.Equal(t, got.Sources[1], markup.ImageProps{URL: "/assets/images/gopher-thumbnail-400x300.webp", IsSource: true, MimeType: markup.ImageMimeTypeWebP, Width: ptr.IntPtr(400), Height: ptr.IntPtr(300)})
		assert.Equal(t, got.Sources[2], markup.ImageProps{URL: "/assets/images/gopher-thumbnail-400x300.jpg", IsSource: true, MimeType: markup.ImageMimeTypeJPG, Width: ptr.IntPtr(400), Height: ptr.IntPtr(300)})
		assert.Equal(t, got.Sources[3], markup.ImageProps{URL: "/assets/images/gopher-mobile-768x1024.avif", IsSource: true, MimeType: markup.ImageMimeTypeAVIF, Width: ptr.IntPtr(768), Height: ptr.IntPtr(1024)})
		assert.Equal(t, got.Sources[4], markup.ImageProps{URL: "/assets/images/gopher-mobile-768x1024.webp", IsSource: true, MimeType: markup.ImageMimeTypeWebP, Width: ptr.IntPtr(768), Height: ptr.IntPtr(1024)})
		assert.Equal(t, got.Sources[5], markup.ImageProps{URL: "/assets/images/gopher-mobile-768x1024.jpg", IsSource: true, MimeType: markup.ImageMimeTypeJPG, Width: ptr.IntPtr(768), Height: ptr.IntPtr(1024)})
		assert.Equal(t, got.Sources[6], markup.ImageProps{URL: "/assets/images/gopher-tablet-1024x1365.avif", IsSource: true, MimeType: markup.ImageMimeTypeAVIF, Width: ptr.IntPtr(1024), Height: ptr.IntPtr(1365)})
		assert.Equal(t, got.Sources[7], markup.ImageProps{URL: "/assets/images/gopher-tablet-1024x1365.webp", IsSource: true, MimeType: markup.ImageMimeTypeWebP, Width: ptr.IntPtr(1024), Height: ptr.IntPtr(1365)})
		assert.Equal(t, got.Sources[8], markup.ImageProps{URL: "/assets/images/gopher-tablet-1024x1365.jpg", IsSource: true, MimeType: markup.ImageMimeTypeJPG, Width: ptr.IntPtr(1024), Height: ptr.IntPtr(1365)})
		assert.Equal(t, got.Sources[9], markup.ImageProps{URL: "/assets/images/gopher-desktop-1440x1920.avif", IsSource: true, MimeType: markup.ImageMimeTypeAVIF, Width: ptr.IntPtr(1440), Height: ptr.IntPtr(1920)})
		assert.Equal(t, got.Sources[10], markup.ImageProps{URL: "/assets/images/gopher-desktop-1440x1920.webp", IsSource: true, MimeType: markup.ImageMimeTypeWebP, Width: ptr.IntPtr(1440), Height: ptr.IntPtr(1920)})
		assert.Equal(t, got.Sources[11], markup.ImageProps{URL: "/assets/images/gopher-desktop-1440x1920.jpg", IsSource: true, MimeType: markup.ImageMimeTypeJPG, Width: ptr.IntPtr(1440), Height: ptr.IntPtr(1920)})
		assert.Equal(t, got.Sources[12], markup.ImageProps{URL: "/assets/images/gopher.avif", IsSource: true, MimeType: markup.ImageMimeTypeAVIF})
		assert.Equal(t, got.Sources[13], markup.ImageProps{URL: "/assets/images/gopher.webp", IsSource: true, MimeType: markup.ImageMimeTypeWebP})
	})

	//t.Log("Errors")
	//{
	//	tt := map[string]struct {
	//		input Image
	//		want  string
	//	}{
	//		"Non existent image": {
	//			input: Image("/assets/images/non_existent_image.jpg"),
	//			want:  "no matches found for glob",
	//		},
	//		"Missing size variants": {
	//			input: Image("/assets/images/missing-sizes.jpg"),
	//			want:  "no matches found for glob",
	//		},
	//	}
	//
	//	for name, test := range tt {
	//		t.Run(name, func(t *testing.T) {
	//			_, err := test.input.PictureMarkup()
	//			assert.Error(t, err)
	//			assert.Contains(t, err.Error(), test.want)
	//		})
	//	}
	//}
}

func TestGetImageProperties(t *testing.T) {
	tt := map[string]struct {
		input   string
		want    imageProperties
		wantErr string
	}{
		"Valid image name": {
			input: "image-desktop-1440x1920.jpg",
			want: imageProperties{
				Width:  1440,
				Height: 1920,
				Mime:   markup.ImageMimeTypeJPG,
			},
		},
		"Invalid image name format": {
			input:   "invalid_format.jpg",
			wantErr: "no regex matches found",
		},
		"Unsupported file extension": {
			input:   "image-desktop-1440x1920.pdf",
			wantErr: "no mime type found for extension",
		},
		"Invalid width": {
			input:   "image-desktop-9wrong9x1920.jpg",
			wantErr: "converting width to int",
		},
		"Invalid height": {
			input:   "image-desktop-1440xinvalid.jpg",
			wantErr: "converting height to int",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := getImageProperties(test.input)
			if test.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}
