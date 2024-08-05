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
	t.Skip()

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
		"Simple Source": {
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
