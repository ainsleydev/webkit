package schemaorg

import "time"

// Article defines an article, such as a news article or piece of investigative report.
// Newspapers and magazines have articles of many different types and this is
// intended to cover them all.
//
// See:
// - https://schema.org/Article
// - https://developers.google.com/search/docs/appearance/structured-data/article
//
//	{
//	   "@context": "https://schema.org",
//	   "@type": "NewsArticle",
//	   "headline": "Title of a News Article",
//	   "image": [
//	       "https://example.com/photos/1x1/photo.jpg",
//	       "https://example.com/photos/4x3/photo.jpg",
//	   ],
//	   "datePublished": "2024-01-05T08:00:00+08:00",
//	   "dateModified": "2024-02-05T09:20:00+08:00",
//	   "author": [
//	       {
//	           "@type": "Person",
//	           "name": "Jane Doe",
//	           "url": "https://example.com/profile/janedoe123"
//	       }
//	   ]
//	}
type Article struct {
	// The title of the article. Consider using a concise title,
	// as long titles may be truncated on some devices.
	Headline string `json:"headline"`

	// The URL to an image that is representative of the article.
	// Use images that are relevant to the article, rather than logos or captions.
	// TODO: Use ImageObject type
	Image []string `json:"image"`

	// The date and time the article was first published, in ISO 8601 format.
	DatePublished time.Time `json:"datePublished"`

	// The date and time the article was most recently modified, in ISO 8601 format.
	DateModified time.Time `json:"dateModified"`

	// The author of the article. To help Google best understand authors
	// across various features, we recommend following the
	Author []Person `json:"author"`

	// The number of words in the text of the Article.
	WordCount int64 `json:"wordCount"`
}

// MarshalJSON implements the json.Marshaler interface to generate
// the JSON-LD for an Article.
func (s *Article) MarshalJSON() ([]byte, error) {
	type Alias Article
	return marshal(&struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		*Alias
	}{
		Context: Context,
		Type:    "NewsArticle",
		Alias:   (*Alias)(s),
	})
}
