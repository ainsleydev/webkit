package markup

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
