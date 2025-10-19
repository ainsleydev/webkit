package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

var media = `
{
   "id": 15,
   "alt": "Alt Text",
   "caption": null,
   "updatedAt": "2024-05-17T18:01:52.169Z",
   "createdAt": "2024-05-17T18:01:52.169Z",
   "url": "https://ainsley.dev/media/image.png",
   "filename": "image.png",
   "mimeType": "image/png",
   "filesize": 743837,
   "width": 1440,
   "height": 4894,
   "sizes": {
       "webp": {
           "url": "https://ainsley.dev/media/image-1440x4894.webp",
           "width": 1440,
           "height": 4894,
           "mimeType": "image/webp",
           "filesize": 136842,
           "filename": "image-1440x4894.webp"
       },
       "avif": {
           "url": "https://ainsley.dev/media/image-1440x4894.avif",
           "width": 1440,
           "height": 4894,
           "mimeType": "image/avif",
           "filesize": 101576,
           "filename": "image-1440x4894.avif"
       },
       "thumbnail": {
           "url": "https://ainsley.dev/media/image-400x300.png",
           "width": 400,
           "height": 300,
           "mimeType": "image/png",
           "filesize": 24434,
           "filename": "image-400x300.png"
       },
       "thumbnail_webp": {
           "url": "https://ainsley.dev/media/image-400x300.webp",
           "width": 400,
           "height": 300,
           "mimeType": "image/webp",
           "filesize": 3856,
           "filename": "image-400x300.webp"
       },
       "thumbnail_avif": {
           "url": "https://ainsley.dev/media/image-400x300.avif",
           "width": 400,
           "height": 300,
           "mimeType": "image/avif",
           "filesize": 4574,
           "filename": "image-400x300.avif"
       },
       "mobile": {
           "url": "https://ainsley.dev/media/image-768x2610.png",
           "width": 768,
           "height": 2610,
           "mimeType": "image/png",
           "filesize": 427862,
           "filename": "image-768x2610.png"
       },
       "mobile_webp": {
           "url": "https://ainsley.dev/media/image-768x2610.webp",
           "width": 768,
           "height": 2610,
           "mimeType": "image/webp",
           "filesize": 55076,
           "filename": "image-768x2610.webp"
       },
       "mobile_avif": {
           "url": "https://ainsley.dev/media/image-768x2610.avif",
           "width": 768,
           "height": 2610,
           "mimeType": "image/avif",
           "filesize": 53918,
           "filename": "image-768x2610.avif"
       },
       "tablet": {
           "url": "https://ainsley.dev/media/image-1024x3480.png",
           "width": 1024,
           "height": 3480,
           "mimeType": "image/png",
           "filesize": 712263,
           "filename": "image-1024x3480.png"
       },
       "tablet_webp": {
           "url": "https://ainsley.dev/media/image-1024x3480.webp",
           "width": 1024,
           "height": 3480,
           "mimeType": "image/webp",
           "filesize": 84742,
           "filename": "image-1024x3480.webp"
       },
       "tablet_avif": {
           "url": "https://ainsley.dev/media/image-1024x3480.avif",
           "width": 1024,
           "height": 3480,
           "mimeType": "image/avif",
           "filesize": 80417,
           "filename": "image-1024x3480.avif"
       }
   }
}`

func TestMedia_UnmarshalJSON(t *testing.T) {
	in := `{
			"id": 15,
				"alt": "Alt Text",
				"caption": null,
				"updatedAt": "2024-05-17T18:01:52.169Z",
				"createdAt": "2024-05-17T18:01:52.169Z",
				"url": "https://ainsley.dev/media/image.png",
				"filename": "image.png",
				"mimeType": "image/png",
				"filesize": 743837,
				"width": 1440,
				"height": 4894,
				"sizes": {
				"webp": {
					"url": "https://ainsley.dev/media/image-1440x4894.webp"
				}
			}
		}`

	tt := map[string]struct {
		input   string
		want    Media
		wantErr bool
	}{
		"OK": {
			input: in,
			want: Media{
				ID:        15,
				CreatedAt: "2024-05-17T18:01:52.169Z",
				UpdatedAt: "2024-05-17T18:01:52.169Z",
				URL:       "https://ainsley.dev/media/image.png",
				Filename:  "image.png",
				MimeType:  "image/png",
				Filesize:  743837,
				Width:     ptr.Float64Ptr(1440),
				Height:    ptr.Float64Ptr(4894),
				Sizes: MediaSizes{
					"webp": MediaSize{
						URL: "https://ainsley.dev/media/image-1440x4894.webp",
					},
				},
				Extra: map[string]any{
					"alt":     "Alt Text",
					"caption": nil,
				},
				RawJSON: []byte(in),
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{wrong}`,
			want:    Media{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var m Media
			err := m.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, m)
		})
	}
}

func TestMedia_ImageMarkup(t *testing.T) {
	t.Parallel()

	var m Media
	err := m.UnmarshalJSON([]byte(media))
	require.NoError(t, err)

	i := m.ImageMarkup()

	assert.Equal(t, i.URL, "https://ainsley.dev/media/image.png")
	assert.Equal(t, "Alt Text", i.Alt)
	assert.Equal(t, 1440, *i.Width)
	assert.Equal(t, 4894, *i.Height)
	assert.Equal(t, "15", i.Attributes["data-payload-media-id"])
	assert.Equal(t, "image.png", i.Attributes["data-payload-media-filename"])
	assert.Equal(t, "726.4 KB", i.Attributes["data-payload-media-filesize"])
}

func TestMediaSize_ImageMarkup(t *testing.T) {
	t.Parallel()

	ms := MediaSize{
		Size:     "mobile",
		URL:      "https://ainsley.dev/media/image-768x2610.png",
		Width:    ptr.Float64Ptr(768),
		Height:   ptr.Float64Ptr(2610),
		MimeType: ptr.StringPtr("image/png"),
		Filesize: ptr.Float64Ptr(427862),
		Filename: ptr.StringPtr("image-768x2610.png"),
	}

	i := ms.ImageMarkup()

	assert.Equal(t, "https://ainsley.dev/media/image-768x2610.png", i.URL)
	assert.Equal(t, 768, *i.Width)
	assert.Equal(t, 2610, *i.Height)
	assert.EqualValues(t, "image/png", i.MimeType)
	assert.False(t, i.IsSource)
	assert.Empty(t, i.Media)
	assert.Empty(t, i.Loading)
	assert.Equal(t, "mobile", i.Attributes["data-payload-size"])
	assert.Equal(t, "417.8 KB", i.Attributes["data-payload-media-filesize"])
	assert.Equal(t, "image-768x2610.png", i.Attributes["data-payload-media-filename"])
}

func TestMedia_ToMarkup(t *testing.T) {
	var m Media
	err := m.UnmarshalJSON([]byte(media))
	require.NoError(t, err)

	p := m.PictureMarkup()

	// Assert main image
	t.Log("Main Image")
	{
		assert.Equal(t, "https://ainsley.dev/media/image.png", p.URL)
		assert.Equal(t, "Alt Text", p.Alt)
		assert.Equal(t, 1440, *p.Width)
		assert.Equal(t, 4894, *p.Height)
		assert.Equal(t, "15", p.Attributes["data-payload-media-id"])
		assert.Equal(t, "image.png", p.Attributes["data-payload-media-filename"])
		assert.Equal(t, "726.4 KB", p.Attributes["data-payload-media-filesize"])
		assert.Len(t, p.Sources, 11)
	}

	t.Log("SourceURL")
	{
		assert.Equal(t, "https://ainsley.dev/media/image-400x300.avif", p.Sources[0].URL)
		assert.EqualValues(t, "image/avif", p.Sources[0].MimeType)
		assert.Equal(t, "thumbnail_avif", p.Sources[0].Name)
		assert.Equal(t, 400, *p.Sources[0].Width)
		assert.Equal(t, 300, *p.Sources[0].Height)
		assert.Equal(t, "thumbnail_avif", p.Sources[0].Attributes["data-payload-size"])
	}
}

func TestMediaSizes_Size(t *testing.T) {
	t.Parallel()

	ms := MediaSizes{
		"size1": {Width: ptr.Float64Ptr(100)},
	}

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		got, err := ms.Size("size1")
		require.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		_, err := ms.Size("wrong")
		assert.Error(t, err)
	})
}

func TestMediaSizes_SortByWidth(t *testing.T) {
	t.Skip()

	tt := map[string]struct {
		input MediaSizes
		want  []string
	}{
		"Sorted widths": {
			input: MediaSizes{
				"size1": {Width: ptr.Float64Ptr(100)},
				"size2": {Width: ptr.Float64Ptr(200)},
				"size3": {Width: ptr.Float64Ptr(300)},
			},
			want: []string{"size1", "size2", "size3"},
		},
		"Nil widths": {
			input: MediaSizes{
				"size1": {Width: ptr.Float64Ptr(200)},
				"size2": {Width: nil},
				"size3": {Width: ptr.Float64Ptr(100)},
				"size4": {Width: nil},
			},
			want: []string{"size3", "size1", "size2", "size4"},
		},
		"Nil widths at end": {
			input: MediaSizes{
				"size1": {Width: ptr.Float64Ptr(200)},
				"size2": {Width: nil},
				"size3": {Width: ptr.Float64Ptr(100)},
				"size4": {Width: ptr.Float64Ptr(300)},
				"size5": {Width: nil},
			},
			want: []string{"size3", "size1", "size4", "size2", "size5"},
		},
		"All nil widths": {
			input: MediaSizes{
				"size1": {Width: nil},
				"size2": {Width: nil},
				"size3": {Width: nil},
			},
			want: []string{"size1", "size2", "size3"},
		},
		"Nil width after non-nil width": {
			input: MediaSizes{
				"size1": {Width: ptr.Float64Ptr(100)},
				"size2": {Width: nil},
			},
			want: []string{"size1", "size2"},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.SortByWidth()
			var sizes []string
			for _, size := range got {
				sizes = append(sizes, size.Size)
			}
			assert.Equal(t, test.want, sizes)
		})
	}
}

func TestMediaFields_Alt(t *testing.T) {
	tt := map[string]struct {
		input        MediaFields
		defaultValue []string
		want         string
	}{
		"Alt Field Present": {
			input:        MediaFields{"alt": "Alt Text"},
			defaultValue: nil,
			want:         "Alt Text",
		},
		"Alt Field Missing With Default": {
			input:        MediaFields{},
			defaultValue: []string{"Default Alt Text"},
			want:         "Default Alt Text",
		},
		"Alt Field Not A String With Default": {
			input:        MediaFields{"alt": 123},
			defaultValue: []string{"Default Alt Text"},
			want:         "Default Alt Text",
		},
		"Alt Field Missing Without Default": {
			input:        MediaFields{},
			defaultValue: nil,
			want:         "",
		},
		"Alt Field Not A String Without Default": {
			input:        MediaFields{"alt": 123},
			defaultValue: nil,
			want:         "",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			m := Media{Extra: test.input}
			assert.Equal(t, test.want, m.Alt(test.defaultValue...))
		})
	}
}

func TestMediaFields_Caption(t *testing.T) {
	tt := map[string]struct {
		input        MediaFields
		defaultValue []string
		want         string
	}{
		"Caption Field Present": {
			input:        MediaFields{"caption": "Caption Text"},
			defaultValue: nil,
			want:         "Caption Text",
		},
		"Caption Field Missing With Default": {
			input:        MediaFields{},
			defaultValue: []string{"Default Caption Text"},
			want:         "Default Caption Text",
		},
		"Caption Field Not A String With Default": {
			input:        MediaFields{"caption": 123},
			defaultValue: []string{"Default Caption Text"},
			want:         "Default Caption Text",
		},
		"Caption Field Missing Without Default": {
			input:        MediaFields{},
			defaultValue: nil,
			want:         "",
		},
		"Caption Field Not A String Without Default": {
			input:        MediaFields{"caption": 123},
			defaultValue: nil,
			want:         "",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			m := Media{Extra: test.input}
			assert.Equal(t, test.want, m.Caption(test.defaultValue...))
		})
	}
}
