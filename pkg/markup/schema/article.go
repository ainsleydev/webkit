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
	Headline string `json:"headline,omitempty"`

	// The URL to an image that is representative of the article.
	// Use images that are relevant to the article, rather than logos or captions.
	// TODO: Use ImageObject type
	Image []string `json:"image,omitempty"`

	// The date and time the article was first published, in ISO 8601 format.
	DatePublished time.Time `json:"datePublished,omitempty"`

	// The date and time the article was most recently modified, in ISO 8601 format.
	DateModified time.Time `json:"dateModified,omitempty"`

	// The author of the article. To help Google best understand authors
	// across various features, we recommend following the
	Author []Person `json:"author,omitempty"`

	// The number of words in the text of the Article.
	// Note: Not currently supported by Google.
	WordCount int64 `json:"wordCount,omitempty"`
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

// https://www.hillwebcreations.com/article-structured-data/#:~:text=But%20if%20the%20page%20is,schema.org%2FBlogPosting).
// It's currently news article, but can be changed to BlogPosting

type TODO struct {
	Context             string    `json:"@context"`
	Type                string    `json:"@type"`
	URL                 string    `json:"url"`
	Headline            string    `json:"headline"`
	Image               string    `json:"image"`
	AlternativeHeadline string    `json:"alternativeHeadline"`
	Description         string    `json:"description"`
	DateCreated         time.Time `json:"dateCreated"`
	DatePublished       time.Time `json:"datePublished"`
	DateModified        time.Time `json:"dateModified"`
	WordCount           string    `json:"wordCount"`
	Keywords            []string  `json:"keywords"`
	InLanguage          string    `json:"inLanguage"`
	IsFamilyFriendly    string    `json:"isFamilyFriendly"`
	Author              struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"author"`
	Creator struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"creator"`
	AccountablePerson struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"accountablePerson"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Logo struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"publisher"`
	CopyrightHolder string `json:"copyrightHolder"`
	CopyrightYear   string `json:"copyrightYear"`
}
