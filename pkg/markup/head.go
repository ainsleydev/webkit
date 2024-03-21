package markup

import "time"

type HeadProps struct {
	Title         string
	Description   string
	Image         string
	DatePublished time.Time
	DateModified  time.Time
}

type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Head outputs the head properties of the document.
func Head() {

}
