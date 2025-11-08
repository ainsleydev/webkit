package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestFindPayloadApps(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		appDef *appdef.Definition
		want   int
	}{
		"Single payload app": {
			appDef: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "cms", Type: appdef.AppTypePayload},
				},
			},
			want: 1,
		},
		"Multiple payload apps": {
			appDef: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "cms", Type: appdef.AppTypePayload},
					{Name: "admin", Type: appdef.AppTypePayload},
				},
			},
			want: 2,
		},
		"Mixed app types": {
			appDef: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "cms", Type: appdef.AppTypePayload},
					{Name: "web", Type: appdef.AppTypeSvelteKit},
					{Name: "api", Type: appdef.AppTypeGoLang},
				},
			},
			want: 1,
		},
		"No payload apps": {
			appDef: &appdef.Definition{
				Apps: []appdef.App{
					{Name: "web", Type: appdef.AppTypeSvelteKit},
				},
			},
			want: 0,
		},
		"Empty apps": {
			appDef: &appdef.Definition{
				Apps: []appdef.App{},
			},
			want: 0,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := findPayloadApps(test.appDef)
			assert.Len(t, got, test.want)
		})
	}
}
