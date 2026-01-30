package pkgjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var payloadMatcher = func(name string) bool {
	return name == "payload" || strings.HasPrefix(name, "@payloadcms/")
}

func TestUpdateDependencies(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg  *PackageJSON
		want int
	}{
		"Updates matching dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
					"react":   "^18.0.0",
				},
			},
			want: 1,
		},
		"Updates across all dependency types": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "3.0.0",
				},
				DevDependencies: map[string]string{
					"@payloadcms/eslint-config": "3.0.0",
				},
				PeerDependencies: map[string]string{
					"@payloadcms/db-postgres": "3.0.0",
				},
			},
			want: 3,
		},
		"No updates when no matches": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"react": "^18.0.0",
				},
			},
			want: 0,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := UpdateDependencies(test.pkg, payloadMatcher, func(name, version string) string {
				return "3.1.0"
			})

			assert.Len(t, result.Updated, test.want)
			if test.want > 0 {
				assert.NotEmpty(t, result.OldVersions)
			}
		})
	}
}

func TestPayloadMatcher(t *testing.T) {
	t.Parallel()

	matcher := func(name string) bool {
		return name == "payload" || strings.HasPrefix(name, "@payloadcms/")
	}

	assert.True(t, matcher("payload"))
	assert.True(t, matcher("@payloadcms/db-postgres"))
	assert.True(t, matcher("@payloadcms/richtext-lexical"))
	assert.False(t, matcher("react"))
	assert.False(t, matcher("payloadjs"))
	assert.False(t, matcher("@types/node"))
}

func TestUpdateDependencies_SkipsDowngrades(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		pkg         *PackageJSON
		wantUpdated int
		wantSkipped int
	}{
		"Skips downgrade in dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "^4.0.0",
				},
			},
			wantUpdated: 0,
			wantSkipped: 1,
		},
		"Allows upgrade in dependencies": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "^2.0.0",
				},
			},
			wantUpdated: 1,
			wantSkipped: 0,
		},
		"Allows same version": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload": "^3.1.0",
				},
			},
			wantUpdated: 1,
			wantSkipped: 0,
		},
		"Skips downgrade in devDependencies": {
			pkg: &PackageJSON{
				DevDependencies: map[string]string{
					"@payloadcms/eslint-config": "5.0.0",
				},
			},
			wantUpdated: 0,
			wantSkipped: 1,
		},
		"Mixed upgrade and downgrade": {
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"payload":                      "^2.0.0",
					"@payloadcms/richtext-lexical": "^4.0.0",
				},
			},
			wantUpdated: 1,
			wantSkipped: 1,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := UpdateDependencies(test.pkg, payloadMatcher, func(name, version string) string {
				return "3.1.0"
			})

			assert.Len(t, result.Updated, test.wantUpdated)
			assert.Len(t, result.Skipped, test.wantSkipped)
		})
	}
}

func TestStripVersionPrefix(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Caret":          {input: "^3.0.0", want: "3.0.0"},
		"Tilde":          {input: "~3.0.0", want: "3.0.0"},
		"Greater equal":  {input: ">=3.0.0", want: "3.0.0"},
		"Less equal":     {input: "<=3.0.0", want: "3.0.0"},
		"Greater":        {input: ">3.0.0", want: "3.0.0"},
		"Less":           {input: "<3.0.0", want: "3.0.0"},
		"Exact":          {input: "3.0.0", want: "3.0.0"},
		"Equals prefix":  {input: "=3.0.0", want: "3.0.0"},
		"With spaces":    {input: " ^3.0.0 ", want: "3.0.0"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, StripVersionPrefix(test.input))
		})
	}
}

func TestIsDowngrade(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		oldVersion string
		newVersion string
		want       bool
	}{
		"Downgrade major":         {oldVersion: "^4.0.0", newVersion: "^3.0.0", want: true},
		"Downgrade minor":         {oldVersion: "^3.2.0", newVersion: "^3.1.0", want: true},
		"Downgrade patch":         {oldVersion: "^3.0.2", newVersion: "^3.0.1", want: true},
		"Upgrade major":           {oldVersion: "^2.0.0", newVersion: "^3.0.0", want: false},
		"Upgrade minor":           {oldVersion: "^3.0.0", newVersion: "^3.1.0", want: false},
		"Upgrade patch":           {oldVersion: "^3.0.0", newVersion: "^3.0.1", want: false},
		"Same version":            {oldVersion: "^3.0.0", newVersion: "^3.0.0", want: false},
		"Different prefixes":      {oldVersion: "~3.1.0", newVersion: "^3.0.0", want: true},
		"Exact to caret upgrade":  {oldVersion: "3.0.0", newVersion: "^4.0.0", want: false},
		"Invalid old version":     {oldVersion: "invalid", newVersion: "^3.0.0", want: false},
		"Invalid new version":     {oldVersion: "^3.0.0", newVersion: "invalid", want: false},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, IsDowngrade(test.oldVersion, test.newVersion))
		})
	}
}

func TestFormatVersion(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		version  string
		useExact bool
		want     string
	}{
		"Exact version": {version: "3.0.0", useExact: true, want: "3.0.0"},
		"Caret version": {version: "3.0.0", useExact: false, want: "^3.0.0"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := FormatVersion(test.version, test.useExact)
			assert.Equal(t, test.want, got)
		})
	}
}
