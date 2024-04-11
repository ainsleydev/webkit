package wordpress

import (
	"fmt"
)

type (
	// Media represents a media item in WordPress.
	Media struct {
		ID           int          `json:"id,omitempty"`
		Date         string       `json:"date,omitempty"`
		DateGMT      string       `json:"date_gmt,omitempty"`
		GUID         GUID         `json:"guid,omitempty"`
		Link         string       `json:"link,omitempty"`
		Modified     string       `json:"modified,omitempty"`
		ModifiedGMT  string       `json:"modifiedGMT,omitempty"`
		Password     string       `json:"password,omitempty"`
		Slug         string       `json:"slug,omitempty"`
		Status       string       `json:"status,omitempty"`
		Type         string       `json:"type,omitempty"`
		Title        Title        `json:"title,omitempty"`
		Author       int          `json:"author,omitempty"`
		MediaStatus  string       `json:"comment_status,omitempty"`
		PingStatus   string       `json:"ping_status,omitempty"`
		AltText      string       `json:"alt_text,omitempty"`
		Caption      any          `json:"caption,omitempty"`
		Description  any          `json:"description,omitempty"`
		MediaType    string       `json:"media_type,omitempty"`
		MediaDetails MediaDetails `json:"media_details,omitempty"`
		Post         int          `json:"post,omitempty"`
		SourceURL    string       `json:"source_url,omitempty"`
	}
	// MediaUploadOptions represents options for uploading media.
	MediaUploadOptions struct {
		Filename    string
		ContentType string
		Data        []byte
	}
	// MediaDetails represents details of a media item.
	MediaDetails struct {
		Raw       string                 `json:"raw,omitempty"`
		Rendered  string                 `json:"rendered,omitempty"`
		Width     int                    `json:"width,omitempty"`
		Height    int                    `json:"height,omitempty"`
		File      string                 `json:"file,omitempty"`
		Sizes     MediaDetailsSizes      `json:"sizes,omitempty"`
		ImageMeta map[string]interface{} `json:"image_meta,omitempty"`
	}
	// MediaDetailsSizes represents different sizes of a media item.
	MediaDetailsSizes struct {
		Thumbnail MediaDetailsSizesItem `json:"thumbnail,omitempty"`
		Medium    MediaDetailsSizesItem `json:"medium,omitempty"`
		Large     MediaDetailsSizesItem `json:"large,omitempty"`
		SiteLogo  MediaDetailsSizesItem `json:"site-logo,omitempty"`
	}
	// MediaDetailsSizesItem represents details of a specific size of a media item.
	MediaDetailsSizesItem struct {
		File      string `json:"file,omitempty"`
		Width     int    `json:"width,omitempty"`
		Height    int    `json:"height,omitempty"`
		MimeType  string `json:"mime_type,omitempty"`
		SourceURL string `json:"source_url,omitempty"`
	}
)

// Media retrieves information about a media item by its ID.
func (c *Client) Media(id int) (Media, error) {
	path := fmt.Sprintf("/media/%d", id)
	var m Media
	return m, c.GetAndUnmarshal(path, &m)
}
