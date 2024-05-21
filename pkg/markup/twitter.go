package markup

// TwitterCard represents web page information according to Twitter Card specifications.
// See https://developer.twitter.com/en/docs/twitter-for-websites/cards/guides/getting-started for more details.
type TwitterCard struct {
	// Required properties
	Site        string `json:"twitter:site,omitempty"`        // The Twitter @username the card should be attributed to.
	Creator     string `json:"twitter:creator,omitempty"`     // The Twitter @username of the content creator.
	Title       string `json:"twitter:title"`                 // A concise title for the related content.
	Description string `json:"twitter:description,omitempty"` // A description that concisely summarizes the content.
	Image       string `json:"twitter:image,omitempty"`       // A URL to a unique image representing the content of the page.

	// Optional properties
	ImageAlt string `json:"twitter:image:alt,omitempty"` // A text description of the image conveying the essential nature of an image to users who are visually impaired.
}
