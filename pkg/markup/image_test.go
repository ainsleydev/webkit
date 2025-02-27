package markup

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type imageProviderMock struct {
	Props ImageProps
}

func (i imageProviderMock) ImageMarkup() ImageProps {
	return i.Props
}

func TestImageProps_Render(t *testing.T) {
	tt := map[string]struct {
		input func() ImageProps
		want  string
	}{
		"Simple SourceURL": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL:      "https://example.com/image.jpg",
						IsSource: true,
						MimeType: ImageMimeTypeWebP,
					},
				})
			},
			want: `<source srcset="https://example.com/image.jpg" type="image/webp" />`,
		},
		"Simple Image with Alt": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL: "https://example.com/image.jpg",
					},
				}, ImageWithAlt("Alternative"))
			},
			want: `<img src="https://example.com/image.jpg" alt="Alternative" />`,
		},
		"Image with Width and Height": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL: "https://example.com/image.jpg",
					},
				}, ImageWithWidth(300), ImageWithHeight(200))
			},
			want: `<img src="https://example.com/image.jpg" width="300" height="200" />`,
		},
		"Image with Loading Loading": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL: "https://example.com/image.jpg",
					},
				}, ImageWithLazyLoading())
			},
			want: `<img src="https://example.com/image.jpg" loading="lazy" />`,
		},
		"Image with Eager Loading": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL: "https://example.com/image.jpg",
					},
				}, ImageWithEagerLoading())
			},
			want: `<img src="https://example.com/image.jpg" loading="eager" />`,
		},
		"Image with Custom Attributes": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL: "https://example.com/image.jpg",
					},
				}, ImageWithAttribute("data-id", "main-image"))
			},
			want: `<img src="https://example.com/image.jpg" data-id="main-image" />`,
		},
		"SourceURL with Media Query": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL:      "https://example.com/image.jpg",
						IsSource: true,
						MimeType: ImageMimeTypeWebP,
						Media:    "(min-width: 600px)",
					},
				})
			},
			want: `<source srcset="https://example.com/image.jpg" type="image/webp" media="(min-width: 600px)" />`,
		},
		"SourceURL with Width and Height": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL:      "https://example.com/image.jpg",
						IsSource: true,
					},
				}, ImageWithWidth(300), ImageWithHeight(200))
			},
			want: `<source srcset="https://example.com/image.jpg" width="300" height="200" />`,
		},
		"SourceURL with Custom Attributes": {
			input: func() ImageProps {
				return Image(&imageProviderMock{
					ImageProps{
						URL:      "https://example.com/image.jpg",
						IsSource: true,
					},
				}, ImageWithAttribute("data-id", "main-source"))
			},
			want: `<source srcset="https://example.com/image.jpg" data-id="main-source" />`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			err := test.input().Render(context.Background(), &buf)
			assert.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}
