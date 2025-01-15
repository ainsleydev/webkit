package payload

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ainsleydev/webkit/pkg/util/ptr"

	"github.com/perimeterx/marshmallow"

	"github.com/ainsleydev/webkit/pkg/markup"
)

// Media defines the fields for media when they are uploaded to PayloadCMS.
//
// See: https://payloadcms.com/docs/upload/overview
type Media struct {
	// The ID of the block, this is generated by Payload and is used to
	// uniquely identify the block.
	ID float64 `json:"id"`

	// Initial media fields, these are also defined in each
	// media size.
	//
	// As per the Payload docs: filename, mimeType, and filesize fields
	// will be automatically added to the upload Collection.
	URL      string   `json:"url"`
	Filename string   `json:"filename"`
	Filesize float64  `json:"filesize"`
	MimeType string   `json:"mimeType"`
	Width    *float64 `json:"width,omitempty"`
	Height   *float64 `json:"height,omitempty"`

	// Key value map of media sizes.
	Sizes MediaSizes `json:"sizes,omitempty"`

	// Timestamps for when the item was created and last updated.
	// These are included by default from Payload.
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`

	// Arbitrary key-value pairs of any other fields that appear within
	// the schema but are not defined in the struct.
	Extra MediaFields `json:"-"`

	// RawJSON is the raw byte slice of the block, which can be used to decode
	// the block into a specific type.
	RawJSON json.RawMessage `json:"-"`
}

// MediaSizes defines a dictionary of media sizes by size name
// (e.g. "small", "medium", "large").
type MediaSizes map[string]MediaSize

// Size returns a MediaSize from the sizes map by key or an
// error if it doesn't exist.
func (ms MediaSizes) Size(name string) (MediaSize, error) {
	size, ok := ms[name]
	if !ok {
		return MediaSize{}, fmt.Errorf("media size not found: %s", name)
	}
	return size, nil
}

// MediaSize defines the fields for the different sizes of media when they
// are uploaded to PayloadCMS.
type MediaSize struct {
	Size     string   `json:"-"` // Name of the media size 	e.g. (thumbnail, small, medium, large)
	URL      string   `json:"url,omitempty"`
	Filename *string  `json:"filename,omitempty"`
	Filesize *float64 `json:"filesize,omitempty"`
	MimeType *string  `json:"mimeType,omitempty"`
	Width    *float64 `json:"width,omitempty"`
	Height   *float64 `json:"height,omitempty"`
}

// ImageMarkup implements the markup.ImageProvider interface and transforms the
// MediaSize item into a markup.ImageProps type ready for rendering an <img>
// to the DOM.
func (ms MediaSize) ImageMarkup() markup.ImageProps {
	// Create attributes map with size information
	attributes := markup.Attributes{
		"data-payload-size": ms.Size,
	}

	// Add optional filesize if present
	if ms.Filesize != nil {
		attributes["data-payload-media-filesize"] = formatFileSize(*ms.Filesize)
	}

	// Add optional filename if present
	if ms.Filename != nil {
		attributes["data-payload-media-filename"] = *ms.Filename
	}

	return markup.ImageProps{
		URL:        ms.URL,
		IsSource:   false,
		MimeType:   markup.ImageMimeType(ptr.String(ms.MimeType)),
		Width:      sizeToIntPointer(ms.Width),
		Height:     sizeToIntPointer(ms.Height),
		Attributes: attributes,
	}
}

// UnmarshalJSON unmarshals the JSON data into the Media type.
// This method is used to extract known fields and assign the remaining
// fields to the fields map.
func (m *Media) UnmarshalJSON(data []byte) error {
	var temp Media
	result, err := marshmallow.Unmarshal(data,
		&temp,
		marshmallow.WithExcludeKnownFieldsFromMap(true),
	)
	if err != nil {
		return err
	}

	*m = temp
	m.RawJSON = data
	m.Extra = result

	for k, v := range m.Sizes {
		m.Sizes[k] = v
	}

	return nil
}

// ImageMarkup implements the markup.ImageProvider interface and transforms the Media item
// into a markup.ImageProps type ready for rendering an <img> to the DOM.
func (m *Media) ImageMarkup() markup.ImageProps {
	return markup.ImageProps{
		URL:      m.URL,
		Alt:      m.Alt(),
		IsSource: false,
		MimeType: markup.ImageMimeType(m.MimeType),
		Loading:  "",
		Width:    sizeToIntPointer(m.Width),
		Height:   sizeToIntPointer(m.Height),
		Attributes: markup.Attributes{
			"data-payload-media-id":       fmt.Sprintf("%v", m.ID),
			"data-payload-media-filename": m.Filename,
			"data-payload-media-filesize": formatFileSize(m.Filesize),
		},
	}
}

// PictureMarkup implements the markup.PictureProvider interface and transforms the Media item
// into a markup.PictureProps type ready for rendering a <picture> the DOM.
func (m *Media) PictureMarkup() markup.PictureProps {
	return markup.PictureProps{
		URL:     m.URL,
		Alt:     m.Alt(),
		Sources: m.Sizes.toMarkup(),
		Width:   sizeToIntPointer(m.Width),
		Height:  sizeToIntPointer(m.Height),
		Attributes: markup.Attributes{
			"data-payload-media-id":       fmt.Sprintf("%v", m.ID),
			"data-payload-media-filename": m.Filename,
			"data-payload-media-filesize": formatFileSize(m.Filesize),
		},
	}
}

// mediaByOrder implements sort.Interface for sorting MediaSize by a predefined order.
type mediaByOrder []MediaSize

// Define the order for the sizes
var sizeOrder = map[string]int{
	"thumbnail_avif": 1,
	"thumbnail_webp": 2,
	"thumbnail":      3,
	"mobile_avif":    4,
	"mobile_webp":    5,
	"mobile":         6,
	"tablet_avif":    7,
	"tablet_webp":    8,
	"tablet":         9,
	"desktop_avif":   10,
	"desktop_webp":   11,
	"desktop":        12,
	"avif":           13,
	"webp":           14,
}

func (a mediaByOrder) Len() int { return len(a) }
func (a mediaByOrder) Less(i, j int) bool {
	// Get the order index for each size
	orderI, okI := sizeOrder[a[i].Size]
	orderJ, okJ := sizeOrder[a[j].Size]
	if !okI {
		orderI = len(sizeOrder) + 1 // Any unknown size should go to the end
	}
	if !okJ {
		orderJ = len(sizeOrder) + 1
	}
	return orderI < orderJ
}
func (a mediaByOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// SortByWidth sorts the media sizes by width from lowest to highest.
// If a width is nil, it will consistently appear at the end.
func (ms MediaSizes) SortByWidth() []MediaSize {
	var sorted mediaByOrder

	for key, mediaSize := range ms {
		mediaSize.Size = key
		sorted = append(sorted, mediaSize)
	}

	// Sort the slice by width
	sort.Sort(sorted)

	// Convert sorted slice back to original format
	result := make([]MediaSize, len(sorted))
	copy(result, sorted)
	return result
}

// toMarkup transforms media sizes into a slice of ImageProps ready for
// rendering onto the DOM.
func (ms MediaSizes) toMarkup() []markup.ImageProps {
	images := make([]markup.ImageProps, len(ms))
	index := 0
	for _, img := range ms.SortByWidth() {
		attr := markup.Attributes{
			"data-payload-size": img.Size,
		}
		if img.Filesize != nil {
			attr["data-payload-media-filesize"] = formatFileSize(*img.Filesize)
		}
		if img.Filename != nil {
			attr["data-payload-media-filename"] = *img.Filename
		}
		images[index] = markup.ImageProps{
			URL:        img.URL,
			IsSource:   true,
			Name:       img.Size,
			Width:      sizeToIntPointer(img.Width),
			Height:     sizeToIntPointer(img.Height),
			MimeType:   markup.ImageMimeType(ptr.String(img.MimeType)),
			Attributes: attr,
		}
		// Ensure media=max-width isn't outputted.
		if img.Size == "webp" || img.Size == "avif" {
			images[index].Width = nil
			images[index].Height = nil
		}
		index++
	}
	return images
}

// MediaFields defines a dictionary of arbitrary fields that are not
// defined in the PayloadCMS schema.
type MediaFields map[string]any

// Alt returns the alt text for the media item if it's defined as
// a field, otherwise it returns the first defaultValue if provided,
// or an empty string if no defaultValue is given.
func (m *Media) Alt(defaultValue ...string) string {
	altText := m.Extra.String("alt")
	if altText == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return altText
}

// Caption returns the caption text for the media item if it's defined as
// a field, otherwise it returns the first defaultValue if provided,
// or an empty string if no defaultValue is given.
func (m *Media) Caption(defaultValue ...string) string {
	captionText := m.Extra.String("caption")
	if captionText == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return captionText
}

// String obtains a string from the key value map.
func (m MediaFields) String(key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func sizeToIntPointer(f *float64) *int {
	if f == nil {
		return nil
	}
	intValue := int(*f)
	return &intValue
}

func formatFileSize(size float64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%.1f B", size)
	}
	if size < unit*unit {
		return fmt.Sprintf("%.1f KB", size/unit)
	}
	return fmt.Sprintf("%.1f MB", size/(unit*unit))
}
