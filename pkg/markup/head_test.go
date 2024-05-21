package markup

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHead(t *testing.T) {
	t.Run("Empty Values", func(t *testing.T) {
		err := HeadTemplate.Execute(&bytes.Buffer{}, nil)
		assert.NoError(t, err)
	})

	t.Run("Simple Title & Description", func(t *testing.T) {
		h := HeadProps{
			Title:       "Hello, World!",
			Description: "This is a test description.",
		}
		buf := bytes.Buffer{}
		err := HeadTemplate.Execute(&buf, h)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), `<title>Hello, World!</title>`)
		assert.Contains(t, buf.String(), `<meta name="description" content="This is a test description." />`)
	})
}
