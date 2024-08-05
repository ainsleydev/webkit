package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNavigation_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		withEnv bool
		want    Navigation
		wantErr bool
	}{
		"OK": {
			input: `{
				"header": [
					{
						"id": "1",
						"title": "Home",
						"url": "/home",
						"children": [],
						"extraField": "someValue"
					}
				],
				"footer": [
					{
						"id": "3",
						"title": "Contact",
						"url": "/contact",
						"children": [],
						"extraField": "someValue"
					}
				],
				"extraTab": "someValue"
			}`,
			want: Navigation{
				Header: NavigationItems{
					{
						ID:       "1",
						Title:    "Home",
						URL:      "/home",
						Children: NavigationItems{},
						Fields: map[string]any{
							"extraField": "someValue",
						},
					},
				},
				Footer: NavigationItems{
					{
						ID:       "3",
						Title:    "Contact",
						URL:      "/contact",
						Children: NavigationItems{},
						Fields: map[string]any{
							"extraField": "someValue",
						},
					},
				},
				Tabs: map[string]interface{}{
					"extraTab": "someValue",
				},
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{header: []}`,
			want:    Navigation{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var n Navigation
			err := n.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, n)
		})
	}
}

func TestNavigationItem_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		withEnv bool
		want    NavigationItem
		wantErr bool
	}{
		"OK": {
			input: `{
				"id": "1",
				"title": "Home",
				"url": "/home",
				"children": [],
				"extraField": "someValue"
			}`,
			want: NavigationItem{
				ID:       "1",
				Title:    "Home",
				URL:      "/home",
				Children: NavigationItems{},
				Fields: map[string]interface{}{
					"extraField": "someValue",
				},
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{id: wrong}`,
			want:    NavigationItem{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var n NavigationItem
			err := n.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, n)
		})
	}
}

func TestNavigationItems_Len(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input NavigationItems
		want  int
	}{
		"Empty": {
			input: NavigationItems{},
			want:  0,
		},
		"With Length": {
			input: NavigationItems{
				{ID: "1"},
				{ID: "2"},
			},
			want: 2,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.input.Len())
		})
	}
}

func TestNavigationItems_Walk(t *testing.T) {
	t.Parallel()

	items := NavigationItems{
		{Title: "Home", URL: "/"},
		{
			Title: "About", URL: "/about",
			Children: NavigationItems{
				{Title: "Team", URL: "/about/team"},
				{Title: "History", URL: "/about/history"},
			},
		},
		{Title: "Contact", URL: "/contact"},
	}

	var visitedItems []string
	walker := func(_ int, item *NavigationItem) {
		visitedItems = append(visitedItems, item.Title+"("+item.URL+")")
	}

	items.Walk(walker)

	want := []string{"Home(/)", "About(/about)", "Team(/about/team)", "History(/about/history)", "Contact(/contact)"}

	assert.Equal(t, want, visitedItems)
}

func TestNavigationItems_MaxDepth(t *testing.T) {
	t.Parallel()
	t.Skip()

	tt := map[string]struct {
		input NavigationItems
		want  int
	}{
		"Empty list": {
			input: NavigationItems{},
			want:  0,
		},
		"Single item with no children": {
			input: NavigationItems{
				{Title: "Home", URL: "/"},
			},
			want: 1,
		},
		"Single item with children": {
			input: NavigationItems{
				{
					Title: "About", URL: "/about",
					Children: NavigationItems{
						{Title: "Team", URL: "/about/team"},
					},
				},
			},
			want: 2,
		},
		"Nested children": {
			input: NavigationItems{
				{Title: "Home", URL: "/"},
				{
					Title: "About", URL: "/about",
					Children: NavigationItems{
						{Title: "Team", URL: "/about/team"},
						{Title: "History", URL: "/about/history"},
					},
				},
				{Title: "Contact", URL: "/contact"},
			},
			want: 3,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.MaxDepth()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNavigationItem_HasChildren(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input NavigationItem
		want  bool
	}{
		"Without Children": {
			input: NavigationItem{},
			want:  false,
		},
		"With Children": {
			input: NavigationItem{
				Children: NavigationItems{
					{Title: "Team", URL: "/about/team"},
				},
			},
			want: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.HasChildren()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNavigationItem_IsRelativeURL(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input NavigationItem
		want  bool
	}{
		"Relative": {
			input: NavigationItem{URL: "/about"},
			want:  true,
		},
		"Absolute": {
			input: NavigationItem{URL: "https://example.com/about"},
			want:  false,
		},
		"Error": {
			input: NavigationItem{URL: ":://wrong"},
			want:  false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.IsRelativeURL()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNavigationItem_IsActive(t *testing.T) {
	t.Parallel()
}
