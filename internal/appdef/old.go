package appdef

import schemaorg "github.com/ainsleydev/webkit/pkg/markup/schema"

type definitionOld struct { //nolint
	Name          string
	Slug          string
	URL           string
	Description   string
	WebKitVersion string
	TitleSuffix   string
	Logo          string
	Adapter       string
	Deployment    string
	Images        []string
	Company       schemaorg.Organisation
	Social        Social
	Menus         map[string][]MenuItem
	Params        map[string]any // Other

	Analytics struct {
		Google    string
		Plausible string
	}
}

type MenuItem struct {
	Identifier string
	Name       string
	URL        string
	Weight     string
	Children   []MenuItem
	Params     map[string]any
}

type Contact struct {
	Email       string
	Phone       string
	Coordinates struct{ Latitude, Longitude string }
	Address     schemaorg.Address
}

type Social struct {
	Twitter   string
	LinkedIn  string
	GitHub    string
	YouTube   string
	Facebook  string
	Instagram string
	TikTok    string
	Dribbble  string
}
