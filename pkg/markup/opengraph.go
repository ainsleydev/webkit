package markup

import "time"

// OpenGraph represents web page information according to OGP.
//
// See https://ogp.me/ for more details.
type OpenGraph struct {
	// Basic Metadata - https://ogp.me/#metadata
	Type  string           `json:"type"`
	Title string           `json:"title"`
	URL   string           `json:"url"`
	Image []OpengraphImage `json:"image"`

	// Optional Metadata - https://ogp.me/#optional
	Audio           []OpengraphAudio `json:"audio,omitempty"`
	Description     string           `json:"description,omitempty"`
	Determiner      string           `json:"determiner,omitempty"`
	Locale          string           `json:"locale,omitempty"`
	LocaleAlt       []string         `json:"locale_alternate,omitempty"`
	Video           []OpengraphVideo `json:"video,omitempty"`
	SiteName        string           `json:"site_name,omitempty"`
	SeeAlso         []string         `json:"see_also,omitempty"`
	FacebookAccount string           `json:"fb:app_id,omitempty"`
}

// OpengraphImage represents a structure of "og:image".
// "og:image" might have following properties:
//
//   - og:image:url
//   - og:image:secure_url
//   - og:image:type
//   - og:image:width
//   - og:image:height
//   - og:image:alt
type OpengraphImage struct {
	URL         string `json:"url,omitempty"`
	SecureURL   string `json:"secure_url,omitempty"`
	ContentType string `json:"type,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Alt         string `json:"alt,omitempty"`
}

// OpengraphVideo represents a structure of "og:video".
// "og:video" might have following properties:
//
//   - og:video:url
//   - og:video:secure_url
//   - og:video:type
//   - og:video:width
//   - og:video:height
type OpengraphVideo struct {
	URL         string `json:"url,omitempty"`
	SecureURL   string `json:"secure_url,omitempty"`
	ContentType string `json:"type,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Duration    int    `json:"duration,omitempty"` // Duration in seconds
}

// OpengraphAudio represents a structure of "og:audio".
// "og:audio" might have the following properties:
//
//   - og:audio:url
//   - og:audio:secure_url
//   - og:audio:type
type OpengraphAudio struct {
	URL         string `json:"url,omitempty"`
	SecureURL   string `json:"secure_url,omitempty"`
	ContentType string `json:"type,omitempty"`
}

// OpenGraphArticle represents a structure of "article" type.
// See https://ogp.me/#type_article for more details.
//
//   - article:published_time
//   - article:modified_time
//   - article:author
//   - article:section
//   - article:tag
type OpenGraphArticle struct {
	PublishedTime time.Time `json:"published_time,omitempty"`
	ModifiedTime  time.Time `json:"modified_time,omitempty"`
	Author        string    `json:"author,omitempty"`  // Full name author
	Section       string    `json:"section,omitempty"` // Usually the URL, i.e. Blog or Insights
	Tags          []string  `json:"tags,omitempty"`    // Associated tags with the article.
}
