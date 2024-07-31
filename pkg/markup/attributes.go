package markup

// Attributes specifies additional attributes for an HTML element as key-value pairs.
// This could be ID, Classes, Data attributes or any other attribute that is
// relevant to the element.
type Attributes map[string]string

const (
	// AttributeClass is the class attribute for an HTML element.
	AttributeClass = "class"
	// AttributeID is the ID attribute for an HTML element.
	AttributeID = "id"
)
