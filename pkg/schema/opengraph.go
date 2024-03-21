package seo

import (
	"bytes"
)

// OpenGraph represents web page information according to OGP.
// See https://ogp.me/ for more details.
type OpenGraph struct {
	// SiteName

	// Basic Metadata - https://ogp.me/#metadata
	Title string           `json:"title"`
	Type  string           `json:"type"`
	Image []OpengraphImage `json:"image"`
	URL   string           `json:"url"`

	// Optional Metadata - https://ogp.me/#optional
	Audio       []OpengraphAudio `json:"audio"`
	Description string           `json:"description"`
	Determiner  string           `json:"determiner"`
	Locale      string           `json:"locale"`
	LocaleAlt   []string         `json:"locale_alternate"`
	SiteName    string           `json:"site_name"`
	Video       []OpengraphVideo `json:"video"`
}

// OpengraphImage represents a structure of "og:image".
// "og:image" might have following properties:
//   - og:image:url
//   - og:image:secure_url
//   - og:image:type
//   - og:image:width
//   - og:image:height
//   - og:image:alt
type OpengraphImage struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"` // Content-Type
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Alt       string `json:"alt"`
}

// OpengraphVideo represents a structure of "og:video".
// "og:video" might have following properties:
//   - og:video:url
//   - og:video:secure_url
//   - og:video:type
//   - og:video:width
//   - og:video:height
type OpengraphVideo struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"` // Content-Type
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	// Duration in seconds
	Duration int `json:"duration"`
}

// OpengraphAudio represents a structure of "og:audio".
// "og:audio" might have the following properties:
//   - og:audio:url
//   - og:audio:secure_url
//   - og:audio:type
type OpengraphAudio struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"` // Content-Type
}

func (og OpenGraph) Render() string {
	b := bytes.Buffer{}

	if og.Title != "" {
		b.WriteString("<meta property=\"og:title\" content=\"")
	}

	ff := ""

	ff.
}
