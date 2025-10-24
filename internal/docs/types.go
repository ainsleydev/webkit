package docs

// DocumentType identifies the type of document to generate.
type DocumentType string

const (
	// DocumentTypeAgents represents the AGENTS.md file for AI coding assistants.
	DocumentTypeAgents DocumentType = "AGENTS.md"

	// Future document types can be added here:
	// DocumentTypeReadme       DocumentType = "README.md"
	// DocumentTypeContributing DocumentType = "CONTRIBUTING.md"
)

// String returns the string representation of the document type.
func (d DocumentType) String() string {
	return string(d)
}

// TemplateName returns the template filename for this document type.
func (d DocumentType) TemplateName() string {
	return string(d)
}

// CustomTemplateName returns the custom template filename (with .tmpl extension).
func (d DocumentType) CustomTemplateName() string {
	return string(d) + ".tmpl"
}
