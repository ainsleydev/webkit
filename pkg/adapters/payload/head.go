package payload

import (
	"context"

	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func (a Adapter) Head(ctx context.Context) (markup.HeadProps, error) {
	settings, err := GetSettings(ctx)
	if err != nil {
		return markup.HeadProps{}, err
	}

	return markup.HeadProps{
		Title:       ptr.String(settings.SiteName),
		Description: "Hey",
	}, nil
}
