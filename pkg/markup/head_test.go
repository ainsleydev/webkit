package markup

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	schemaorg "github.com/ainsleydev/webkit/pkg/markup/schema"
)

var props = HeadProps{
	Title:         "Title",
	Description:   "AgencyDescription",
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
		Description: "AgencyDescription",
		Image:       "https://example.com/simage.png",
	},
	Organisation: &schemaorg.Organisation{
		ID:          "https://example.com",
		URL:         "https://example.com",
		LegalName:   "Example Inc.",
		Description: "AgencyDescription",
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
	//			AgencyDescription: "Homepage of the website",
	//			URL:         "https://example.com",
	//		},
	//		{
	//			Type:        "ListItem",
	//			Position:    2,
	//			Name:        "About Us",
	//			AgencyDescription: "About us page",
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

	t.Run("Simple Title & AgencyDescription", func(t *testing.T) {
		p := HeadProps{
			Title:       "Hello, World!",
			Description: "This is a test description.",
			Private:     true,
		}
		buf := bytes.Buffer{}
		err := p.Render(context.TODO(), &buf)

		assert.NoError(t, err)

		doc, err := goquery.NewDocumentFromReader(&buf)

		require.NoError(t, err)
		assert.Contains(t, doc.Find("title").Text(), "Hello, World!")
		description, ok := doc.Find(`meta[name="description"]`).Attr("content")
		assert.True(t, ok)
		assert.Equal(t, "This is a test description.", description)
	})

	t.Run("Full Head", func(t *testing.T) {
		buf := bytes.Buffer{}
		ctx := context.Background()
		ctx = webkitctx.WithHeadSnippet(ctx, "Test", "<script>alert('Hello, World!');</script>")

		err := props.Render(ctx, &buf)
		assert.NoError(t, err)
	})
}
