package markup

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

type pictureProviderMock struct {
	Props PictureProps
}

func (p pictureProviderMock) PictureMarkup() PictureProps {
	return p.Props
}

func TestPictureProps_Render(t *testing.T) {
	tt := map[string]struct {
		input func() PictureProps
		want  string
	}{
		"Image Only": {
			input: func() PictureProps {
				return Picture(&pictureProviderMock{
					PictureProps{
						URL: "https://example.com/image.jpg",
					},
				}, PictureWithAlt("Alternative"))
			},
			want: "<picture>\n  <img src=\"https://example.com/image.jpg\" alt=\"Alternative\" />\n</picture>",
		},
		"Simple SourceURL": {
			input: func() PictureProps {
				return Picture(&pictureProviderMock{
					PictureProps{
						URL: "https://example.com/image.jpg",
						Sources: []ImageProps{
							{
								URL:      "https://example.com/image-thumbnail.avif",
								IsSource: true,
								Width:    ptr.IntPtr(400),
								Media:    "(max-width: 450px)",
								MimeType: "image/avif",
							},
						},
					},
				}, PictureWithAlt("Alternative"))
			},
		},
		"Image with Loading Loading": {
			input: func() PictureProps {
				return Picture(&pictureProviderMock{
					PictureProps{
						URL: "https://example.com/image.jpg",
					},
				}, PictureWithLazyLoading())
			},
			want: `<img src="https://example.com/image.jpg" loading="lazy" />`,
		},
		"Image with Eager Loading": {
			input: func() PictureProps {
				return Picture(&pictureProviderMock{
					PictureProps{
						URL: "https://example.com/image.jpg",
					},
				}, PictureWithEagerLoading())
			},
			want: `<img src="https://example.com/image.jpg" loading="eager" />`,
		},
		"Image with Attribute": {
			input: func() PictureProps {
				return Picture(&pictureProviderMock{
					PictureProps{
						URL: "https://example.com/image.jpg",
					},
				}, PictureWithAttribute("data-id", "main-source"))
			},
			want: `<picture data-id="main-source">`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			err := test.input().Render(context.Background(), &buf)
			assert.NoError(t, err)
			assert.Contains(t, buf.String(), test.want)
		})
	}
}

func TestPictureWithSize(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		name  string
		input PictureProps
		sizes []string
		want  PictureProps
	}{
		"No sizes specified": {
			input: PictureProps{
				URL: "https://example.com/original.jpg",
				Sources: []ImageProps{
					{Name: "thumbnail", URL: "https://example.com/thumbnail.jpg"},
					{Name: "medium", URL: "https://example.com/medium.jpg"},
				},
			},
			sizes: []string{},
			want: PictureProps{
				URL: "https://example.com/original.jpg",
				Sources: []ImageProps{
					{Name: "thumbnail", URL: "https://example.com/thumbnail.jpg"},
					{Name: "medium", URL: "https://example.com/medium.jpg"},
				},
			},
		},
		"Single size match": {
			input: PictureProps{
				Sources: []ImageProps{
					{Name: "thumbnail", URL: "https://example.com/thumbnail.jpg", Width: ptr.IntPtr(100), Height: ptr.IntPtr(100)},
					{Name: "medium", URL: "https://example.com/medium.jpg", Width: ptr.IntPtr(300), Height: ptr.IntPtr(300)},
					{Name: "large", URL: "https://example.com/large.jpg", Width: ptr.IntPtr(600), Height: ptr.IntPtr(600)},
					{Name: "medium-large", URL: "https://example.com/medium-large.jpg", Width: ptr.IntPtr(450), Height: ptr.IntPtr(450)},
				},
			},
			sizes: []string{"medium"},
			want: PictureProps{
				URL:    "https://example.com/medium.jpg",
				Width:  ptr.IntPtr(300),
				Height: ptr.IntPtr(300),
				Sources: []ImageProps{
					{Name: "medium-large", URL: "https://example.com/medium-large.jpg", Width: ptr.IntPtr(450), Height: ptr.IntPtr(450)},
				},
			},
		},
		"Multiple sizes with partial match": {
			input: PictureProps{
				Sources: []ImageProps{
					{Name: "thumbnail", URL: "https://example.com/thumbnail.jpg", Width: ptr.IntPtr(100), Height: ptr.IntPtr(100)},
					{Name: "medium", URL: "https://example.com/medium.jpg", Width: ptr.IntPtr(300), Height: ptr.IntPtr(300)},
					{Name: "large", URL: "https://example.com/large.jpg", Width: ptr.IntPtr(600), Height: ptr.IntPtr(600)},
					{Name: "medium-special", URL: "https://example.com/medium-special.jpg"},
				},
			},
			sizes: []string{"thumbnail", "medium"},
			want: PictureProps{
				URL:    "https://example.com/medium.jpg",
				Width:  ptr.IntPtr(300),
				Height: ptr.IntPtr(300),
				Sources: []ImageProps{
					{Name: "thumbnail", URL: "https://example.com/thumbnail.jpg", Width: ptr.IntPtr(100), Height: ptr.IntPtr(100)},
					{Name: "medium-special", URL: "https://example.com/medium-special.jpg"},
				},
			},
		},
		"No matching sizes": {
			input: PictureProps{
				URL: "https://example.com/original.jpg",
				Sources: []ImageProps{
					{Name: "small", URL: "https://example.com/small.jpg"},
					{Name: "large", URL: "https://example.com/large.jpg"},
				},
			},
			sizes: []string{"nonexistent"},
			want: PictureProps{
				URL: "https://example.com/original.jpg",
				Sources: []ImageProps{
					{Name: "small", URL: "https://example.com/small.jpg"},
					{Name: "large", URL: "https://example.com/large.jpg"},
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			option := PictureWithSize(test.sizes...)
			option(&test.input)

			assert.Equal(t, test.want.URL, test.input.URL, "URLs should match")
			assert.Equal(t, test.want.Width, test.input.Width, "Widths should match")
			assert.Equal(t, test.want.Height, test.input.Height, "Heights should match")

			assert.Equal(t, len(test.want.Sources), len(test.input.Sources), "Number of sources should match")
			for i := range test.want.Sources {
				assert.Equal(t, test.want.Sources[i], test.input.Sources[i], "Source %d should match", i)
			}
		})
	}
}
