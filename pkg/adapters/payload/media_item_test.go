package payload

import (
	"bytes"
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
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
   "url": "/media/image.png",
   "filename": "image.png",
   "mimeType": "image/png",
   "filesize": 743837,
   "width": 1440,
   "height": 4894,
   "sizes": {
       "webp": {
           "url": "/media/image-1440x4894.webp",
           "width": 1440,
           "height": 4894,
           "mimeType": "image/webp",
           "filesize": 136842,
           "filename": "image-1440x4894.webp"
       },
       "avif": {
           "url": "/media/image-1440x4894.avif",
           "width": 1440,
           "height": 4894,
           "mimeType": "image/avif",
           "filesize": 101576,
           "filename": "image-1440x4894.avif"
       },
       "thumbnail": {
           "url": "/media/image-400x300.png",
           "width": 400,
           "height": 300,
           "mimeType": "image/png",
           "filesize": 24434,
           "filename": "image-400x300.png"
       },
       "thumbnail_webp": {
           "url": "/media/image-400x300.webp",
           "width": 400,
           "height": 300,
           "mimeType": "image/webp",
           "filesize": 3856,
           "filename": "image-400x300.webp"
       },
       "thumbnail_avif": {
           "url": "/media/image-400x300.avif",
           "width": 400,
           "height": 300,
           "mimeType": "image/avif",
           "filesize": 4574,
           "filename": "image-400x300.avif"
       },
       "mobile": {
           "url": "/media/image-768x2610.png",
           "width": 768,
           "height": 2610,
           "mimeType": "image/png",
           "filesize": 427862,
           "filename": "image-768x2610.png"
       },
       "mobile_webp": {
           "url": "/media/image-768x2610.webp",
           "width": 768,
           "height": 2610,
           "mimeType": "image/webp",
           "filesize": 55076,
           "filename": "image-768x2610.webp"
       },
       "mobile_avif": {
           "url": "/media/image-768x2610.avif",
           "width": 768,
           "height": 2610,
           "mimeType": "image/avif",
           "filesize": 53918,
           "filename": "image-768x2610.avif"
       },
       "tablet": {
           "url": "/media/image-1024x3480.png",
           "width": 1024,
           "height": 3480,
           "mimeType": "image/png",
           "filesize": 712263,
           "filename": "image-1024x3480.png"
       },
       "tablet_webp": {
           "url": "/media/image-1024x3480.webp",
           "width": 1024,
           "height": 3480,
           "mimeType": "image/webp",
           "filesize": 84742,
           "filename": "image-1024x3480.webp"
       },
       "tablet_avif": {
           "url": "/media/image-1024x3480.avif",
           "width": 1024,
           "height": 3480,
           "mimeType": "image/avif",
           "filesize": 80417,
           "filename": "image-1024x3480.avif"
       }
   }
}`

func TestMedia_Render(t *testing.T) {
	t.Setenv(envPayloadURL, "https://example.com")

	var m Media
	err := m.UnmarshalJSON([]byte(media))
	require.NoError(t, err)

	r := &bytes.Buffer{}
	err = m.Render(context.Background(), r)
	require.NoError(t, err)

	doc, err := goquery.NewDocumentFromReader(r)
	require.NoError(t, err)

	sources := []string{
		`source[srcset="https://example.com/media/image-400x300.png"][media="(max-width: 450px)"][type="image/png"][width="400"][height="300"][data-payload-size="thumbnail"]`,
		`source[srcset="https://example.com/media/image-400x300.avif"][media="(max-width: 450px)"][type="image/avif"][width="400"][height="300"][data-payload-size="thumbnail_avif"]`,
		`source[srcset="https://example.com/media/image-400x300.webp"][media="(max-width: 450px)"][type="image/webp"][width="400"][height="300"][data-payload-size="thumbnail_webp"]`,
		`source[srcset="https://example.com/media/image-768x2610.png"][media="(max-width: 818px)"][type="image/png"][width="768"][height="2610"][data-payload-size="mobile"]`,
		`source[srcset="https://example.com/media/image-768x2610.avif"][media="(max-width: 818px)"][type="image/avif"][width="768"][height="2610"][data-payload-size="mobile_avif"]`,
		`source[srcset="https://example.com/media/image-768x2610.webp"][media="(max-width: 818px)"][type="image/webp"][width="768"][height="2610"][data-payload-size="mobile_webp"]`,
		`source[srcset="https://example.com/media/image-1024x3480.png"][media="(max-width: 1074px)"][type="image/png"][width="1024"][height="3480"][data-payload-size="tablet"]`,
		`source[srcset="https://example.com/media/image-1024x3480.avif"][media="(max-width: 1074px)"][type="image/avif"][width="1024"][height="3480"][data-payload-size="tablet_avif"]`,
		`source[srcset="https://example.com/media/image-1024x3480.webp"][media="(max-width: 1074px)"][type="image/webp"][width="1024"][height="3480"][data-payload-size="tablet_webp"]`,
	}

	for _, src := range sources {
		assert.Equal(t, 1, doc.Find(src).Length())
	}
}

func TestMedia_UnmarshalJSON(t *testing.T) {
	var (
		payloadURL = "https://example.com"
		in         = `{
			"id": 15,
				"alt": "Alt Text",
				"caption": null,
				"updatedAt": "2024-05-17T18:01:52.169Z",
				"createdAt": "2024-05-17T18:01:52.169Z",
				"url": "/media/image.png",
				"filename": "image.png",
				"mimeType": "image/png",
				"filesize": 743837,
				"width": 1440,
				"height": 4894,
				"sizes": {
				"webp": {
					"url": "/media/image-1440x4894.webp"
				}
			}
		}`
	)

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
				URL:       "https://example.com/media/image.png",
				Filename:  "image.png",
				MimeType:  "image/png",
				Filesize:  743837,
				Width:     ptr.Float64Ptr(1440),
				Height:    ptr.Float64Ptr(4894),
				Sizes: MediaSizes{
					"webp": MediaSize{
						URL: "https://example.com/media/image-1440x4894.webp",
					},
				},
				Fields: map[string]any{
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
			t.Setenv(envPayloadURL, payloadURL)
			var m Media
			err := m.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, m)
		})
	}

	t.Run("No Env", func(t *testing.T) {
		var m Media
		err := m.UnmarshalJSON([]byte(in))
		assert.Error(t, err)
	})
}

func TestMediaSizes_SortByWidth(t *testing.T) {
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
	ttt := map[string]struct {
		input MediaFields
		want  string
	}{
		"Alt Field Present": {
			input: MediaFields{"alt": "Alt Text"},
			want:  "Alt Text",
		},
		"Alt Field Missing": {
			input: MediaFields{},
			want:  "",
		},
		"Alt Field Not A String": {
			input: MediaFields{"alt": 123},
			want:  "",
		},
	}

	for name, test := range ttt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, test.input.Alt())
		})
	}
}

func TestMediaFields_Caption(t *testing.T) {
	ttt := map[string]struct {
		input MediaFields
		want  string
	}{
		"Caption Field Present": {
			input: MediaFields{"caption": "Caption Text"},
			want:  "Caption Text",
		},
		"Caption Field Missing": {
			input: MediaFields{},
			want:  "",
		},
		"Caption Field Not A String": {
			input: MediaFields{"caption": 123},
			want:  "",
		},
	}

	for name, test := range ttt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, test.input.Caption())
		})
	}
}
