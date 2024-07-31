package markup

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestImageProps_Render(t *testing.T) {
	tt := map[string]struct {
		input ImageProps
		want  string
	}{
		"Simple Image": {
			input: ImageProps{
				URL: "https://example.com/image.jpg",
				Attributes: Attributes{
					"alt": "An example image",
				},
			},
			want: `<img src="https://example.com/image.jpg" alt="An example image">`,
		},
		"Source Element With Media Query": {
			input: ImageProps{
				URL:      "https://example.com/image.webp",
				IsSource: true,
				Media:    "(min-width: 600px)",
				MimeType: ptr.StringPtr("image/webp"),
			},
			want: `<source srcset="https://example.com/image.webp" type="image/webp" media="(min-width: 600px)" />`,
		},
		"Img Element With Width And Height": {
			input: ImageProps{
				URL:    "https://example.com/image.png",
				Width:  ptr.IntPtr(300),
				Height: ptr.IntPtr(200),
			},
			want: `<img src="https://example.com/image.png" width="300" height="200">`,
		},
		"Source Element With Attributes": {
			input: ImageProps{
				URL:      "https://example.com/image.avif",
				IsSource: true,
				MimeType: ptr.StringPtr("image/avif"),
				Attributes: Attributes{
					"class": "responsive",
				},
			},
			want: `<source srcset="https://example.com/image.avif" type="image/avif" class="responsive" />`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			err := test.input.Render(context.Background(), &buf)
			assert.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}
