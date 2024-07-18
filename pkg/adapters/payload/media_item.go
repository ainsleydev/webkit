package payload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/goccy/go-json"
	"github.com/perimeterx/marshmallow"

	"github.com/ainsleydev/webkit/pkg/adapters/payload/internal/tpl"
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

	// Timestamps for when the media was created and last updated.
	// These are included by default from Payload.
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`

	// Arbitrary key-value pairs of any other fields that appear within
	// the schema but are not defined in the struct.
	Fields MediaFields `json:"-"`

	// RawJSON is the raw byte slice of the block, which can be used to decode
	// the block into a specific type.
	RawJSON json.RawMessage `json:"-"`
}

// MediaSizes defines a dictionary of media sizes by size name
// (e.g. "small", "medium", "large").
type MediaSizes map[string]MediaSize

// MediaSize defines the fields for the different sizes of media when they
// are uploaded to PayloadCMS.
type MediaSize struct {
	Size      string   `json:"-"` // Name of the media size 	e.g. (thumbnail, small, medium, large)
	URL       string   `json:"url,omitempty"`
	Filename  *string  `json:"filename,omitempty"`
	Filesize  *float64 `json:"filesize,omitempty"`
	MimeType  *string  `json:"mimeType,omitempty"`
	Width     *float64 `json:"width,omitempty"`
	Height    *float64 `json:"height,omitempty"`
	MediaAttr string   `json:"media,omitempty"`
}

// Render renders the media block to the provided writer as a
// picture element.
//
// Note: It does not include the <picture> element and it
// expects the caller to wrap the output.
func (m *Media) Render(_ context.Context, w io.Writer) error {
	return tpl.Templates.ExecuteTemplate(w, "picture.html", m)
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

	url := os.Getenv(envPayloadURL)
	if url == "" {
		return errors.New("env var: " + envPayloadURL + " is not set")
	}

	*m = temp
	m.RawJSON = data
	m.Fields = result
	m.URL = url + m.URL

	for k, v := range m.Sizes {
		v.URL = url + v.URL
		m.Sizes[k] = v
	}

	return nil
}

// mediaByWidth implements sort.Interface for sorting MediaSize by Width.
type mediaByWidth []MediaSize

func (a mediaByWidth) Len() int { return len(a) }
func (a mediaByWidth) Less(i, j int) bool {
	// Handle nil width consistently
	if a[i].Width == nil && a[j].Width != nil {
		return false
	} else if a[i].Width != nil && a[j].Width == nil {
		return true
	}

	// Sort by width, then by key for stability
	if a[i].Width == nil || a[j].Width == nil || *a[i].Width == *a[j].Width {
		return a[i].Size < a[j].Size
	}
	return *a[i].Width < *a[j].Width
}
func (a mediaByWidth) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// SortByWidth sorts the media sizes by width from lowest to highest.
// If a width is nil, it will consistently appear at the end.
func (ms MediaSizes) SortByWidth() []MediaSize {
	// Convert map to slice for deterministic sorting
	sorted := make(mediaByWidth, 0, len(ms))
	for key, mediaSize := range ms {
		mediaSize.Size = key
		if mediaSize.Width != nil {
			mediaSize.MediaAttr = fmt.Sprintf("(max-width: %vpx)", *mediaSize.Width+50)
		}
		sorted = append(sorted, mediaSize)
	}
	sort.Sort(sorted)

	// Convert sorted slice back to original format
	result := make([]MediaSize, len(sorted))
	for i, m := range sorted {
		result[i] = m
	}
	return result
}

// MediaFields defines a dictionary of arbitrary fields that are not
// defined in the PayloadCMS schema.
type MediaFields map[string]any

// Alt returns the alt text for the media item if it's defined as
// a field, otherwise an empty string.
func (m MediaFields) Alt() string {
	return m.string("alt")
}

// Caption returns the caption text for the media item if it's defined as
// a field, otherwise an empty string.
func (m MediaFields) Caption() string {
	return m.string("caption")
}

func (m MediaFields) string(key string) string {
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
