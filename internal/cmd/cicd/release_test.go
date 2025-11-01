package cicd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestReleaseWorkflow_FilterApps(t *testing.T) {
	t.Parallel()

	falseVal := false
	trueVal := true

	tt := map[string]struct {
		apps []appdef.App
		want int
	}{
		"All apps with Dockerfile and release enabled": {
			apps: []appdef.App{
				{Name: "cms", Build: appdef.Build{Dockerfile: "Dockerfile", Release: nil}},
				{Name: "web", Build: appdef.Build{Dockerfile: "Dockerfile", Release: &trueVal}},
				{Name: "api", Build: appdef.Build{Dockerfile: "Dockerfile"}},
			},
			want: 3,
		},
		"One app with release disabled": {
			apps: []appdef.App{
				{Name: "cms", Build: appdef.Build{Dockerfile: "Dockerfile"}},
				{Name: "web", Build: appdef.Build{Dockerfile: "Dockerfile", Release: &falseVal}},
				{Name: "api", Build: appdef.Build{Dockerfile: "Dockerfile"}},
			},
			want: 2,
		},
		"No Dockerfile defined": {
			apps: []appdef.App{
				{Name: "cms", Build: appdef.Build{Dockerfile: ""}},
			},
			want: 0,
		},
		"Empty apps list": {
			apps: []appdef.App{},
			want: 0,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var appsToRelease []appdef.App
			for _, app := range test.apps {
				if app.Build.Dockerfile != "" && app.ShouldRelease() {
					appsToRelease = append(appsToRelease, app)
				}
			}

			assert.Len(t, appsToRelease, test.want)
		})
	}
}
