package markup

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	schemaorg "github.com/ainsleydev/webkit/pkg/markup/schema"
)

var props = HeadProps{
	Title:         "Title",
	Description:   "Description",
	Image:         "https://example.com/image.png",
	DatePublished: time.Now(),
	DateModified:  time.Now(),
	Locale:        "en_US",
	URL:           "https://example.com/page",
	IsPage:        true,
	Private:       true,
	Canonical:     "https://example.com/page",
	OpenGraph: &OpenGraph{
		Type:  "website",
		Title: "Sample OpenGraph Title",
		URL:   "https://example.com/page",
		Image: []OpengraphImage{
			{
				URL:         "https://example.com/image-og.png",
				ContentType: "image/png",
				Width:       1200,
				Height:      630,
				Alt:         "Alt",
			},
		},
	},
	Twitter: &TwitterCard{
		Site:        "@example",
		Creator:     "@author",
		Title:       "Title",
		Description: "Description",
		Image:       "https://example.com/simage.png",
	},
	Organisation: &schemaorg.Organisation{
		ID:          "https://example.com",
		URL:         "https://example.com",
		LegalName:   "Example Inc.",
		Description: "Description",
		Logo:        "https://example.com/logo.png",
		SameAs: []string{
			"https://twitter.com/example",
			"https://linkedin.com/company/example",
		},
		Address: schemaorg.Address{
			StreetAddress: "123 Sample St",
			Locality:      "Sample City",
			Region:        "Sample State",
			Country:       "UK",
			PostalCode:    "12345",
		},
	},
	//Navigation: &SchemaOrgNavItemList{
	//	Context: "https://schema.org",
	//	Type:    "ItemList",
	//	ItemListElement: []SchemaOrgItemListElement{
	//		{
	//			Type:        "ListItem",
	//			Position:    1,
	//			Name:        "Home",
	//			Description: "Homepage of the website",
	//			URL:         "https://example.com",
	//		},
	//		{
	//			Type:        "ListItem",
	//			Position:    2,
	//			Name:        "About Us",
	//			Description: "About us page",
	//			URL:         "https://example.com/about",
	//		},
	//	},
	//},
	Other: "<script>alert('Hello, World!');</script>",
}

func TestHead(t *testing.T) {
	t.Run("Empty Values", func(t *testing.T) {
		p := HeadProps{}
		err := p.Render(context.TODO(), &bytes.Buffer{})
		assert.NoError(t, err)
	})

	t.Run("Simple Title & Description", func(t *testing.T) {
		p := HeadProps{
			Title:       "Hello, World!",
			Description: "This is a test description.",
			Private:     true,
		}
		buf := bytes.Buffer{}
		err := p.Render(context.TODO(), &buf)
		fmt.Println(buf.String())

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), `<title>Hello, World!</title>`)
		assert.Contains(t, buf.String(), `<meta name="description" content="This is a test description." />`)
	})

	t.Run("Full Head", func(t *testing.T) {
		buf := bytes.Buffer{}
		ctx := context.Background()
		ctx = webkitctx.WithHeadSnippet(ctx, "Test", "<script>alert('Hello, World!');</script>")

		err := props.Render(ctx, &buf)
		assert.NoError(t, err)
		fmt.Println(buf.String())
	})
}
