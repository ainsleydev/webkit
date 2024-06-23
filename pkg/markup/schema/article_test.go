package schemaorg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArticle_MarshalJSON(t *testing.T) {
	article := Article{
		Headline: "Headline",
		Image: []string{
			"https://example.com/photos/1x1/photo.jpg",
			"https://example.com/photos/4x3/photo.jpg",
		},
		DatePublished: time.Date(2024, 1, 5, 8, 0, 0, 0, time.FixedZone("CST", 8*60*60)),
		DateModified:  time.Date(2024, 2, 5, 9, 20, 0, 0, time.FixedZone("CST", 8*60*60)),
		Author: []Person{
			{
				Name: "Jane Doe",
				URL:  "https://example.com/profile/janedoe123",
			},
		},
		WordCount: 500,
	}

	want := `{
		"@context": "https://schema.org",
		"@type": "NewsArticle",
		"headline": "Headline",
		"image": [
			"https://example.com/photos/1x1/photo.jpg",
			"https://example.com/photos/4x3/photo.jpg"
		],
		"datePublished": "2024-01-05T08:00:00+08:00",
		"dateModified": "2024-02-05T09:20:00+08:00",
		"author": [
			{
				"@type": "Person",
				"name": "Jane Doe",
				"url": "https://example.com/profile/janedoe123"
			}
		],
		"wordCount": 500
	}`

	got, err := article.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, want, string(got))
}
