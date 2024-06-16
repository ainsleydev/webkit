package payload

import (
	"fmt"
	"sort"

	"github.com/perimeterx/marshmallow"
)

// Media defines the fields for media when they are uploaded to PayloadCMS.
type Media struct {
	Id        float64    `json:"id"`
	MediaSize            // The original
	Alt       string     `json:"alt,omitempty"`
	Sizes     MediaSizes `json:"sizes,omitempty"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`

	// Key-value pairs of the media's fields, these pairs are defined by the block's
	// schema and vary depending on the block type.
	Fields map[string]any `json:"-"`
}

// MediaSizes defines a dictionary of media sizes by size name
// (e.g. "small", "medium", "large").
type MediaSizes map[string]MediaSize

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

// UnmarshalJSON unmarshals the JSON data into the Media type.
// This method is used to extract known fields and assign the remaining
// fields to the fields map.
func (m *Media) UnmarshalJSON(data []byte) error {
	var temp Media
	result, err := marshmallow.Unmarshal(data, &temp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*m = temp
	m.Fields = result
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
