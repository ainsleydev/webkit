package static

import (
	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/stretchr/testify/assert"
	"testing"
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
	t.Skip()
}

func TestImage_RootPath(t *testing.T) {

}
