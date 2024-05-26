package payload

import (
	"context"

	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func (a Adapter) Head(ctx context.Context) markup.HeadProps {
	settings := getSettings(ctx)

	if settings == nil {
		return markup.HeadProps{}
	}

	return markup.HeadProps{
		Title: ptr.String(settings.SiteName),
	}
}
