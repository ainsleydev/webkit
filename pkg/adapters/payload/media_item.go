package payload

import (
	"fmt"
	"html/template"
	"io"
	"sort"

	"github.com/ainsleydev/webkit/pkg/adapters/payload/internal/templates"
)

// Media defines the fields for media when they are uploaded to PayloadCMS.
type Media struct {
	Id  float64 `json:"id"`
	Alt string  `json:"alt"`
	// TODO, need to parse RichText
	Caption   []map[string]any `json:"caption,omitempty"`
	MediaSize                  // The original
	Sizes     MediaSizes       `json:"sizes,omitempty"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt" `
}

// MediaSizes defines a dictionary of media sizes by size name
// (e.g. "small", "medium", "large").
type MediaSizes map[string]MediaSize

// MediaSize defines the fields for the different sizes of media when they
// are uploaded to PayloadCMS.
type MediaSize struct {
	Size     string   `json:"-"` // Name of the media size 	e.g. (thumbnail, small, medium, large)
	MediaSrc string   `json:"-"` // The media attribute 	e.g. (min-width: 600px)
	URL      *string  `json:"url,omitempty"`
	Filename *string  `json:"filename,omitempty"`
	Filesize *float64 `json:"filesize,omitempty"`
	MimeType *string  `json:"mimeType,omitempty"`
	Width    *float64 `json:"width,omitempty"`
	Height   *float64 `json:"height,omitempty"`
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
		mediaSize.MediaSrc = fmt.Sprintf("(min-width: %vpx)", *mediaSize.Width+50)
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

// TODO: - Caption on Media
func (m Media) Render(w io.Writer) error {
	tpl, err := template.New("").Parse(templates.Picture)
	if err != nil {
		return err
	}
	return tpl.Execute(w, m)
}
