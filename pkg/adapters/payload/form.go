package payload

import (
	"context"
	"io"

	"github.com/ainsleydev/webkit/pkg/adapters/payload/internal/tpl"
)

// Form defines a singular form collection type in the Form Builder Plugin
// within Payload CMS
// See: https://payloadcms.com/docs/plugins/form-builder
//
// Blocks example within the frontend:
// https://github.com/payloadcms/payload/tree/main/examples/form-builder/next-pages/components/Blocks
//
// TypeScript Bindings:
// https://github.com/payloadcms/payload/blob/main/packages/plugin-form-builder/src/types.ts
type Form struct {
	ID                int                  `json:"id"`
	Title             string               `json:"title"`
	Fields            []FormField          `json:"fields,omitempty"`
	SubmitButtonLabel *string              `json:"submitButtonLabel,omitempty"`
	ConfirmationType  FormConfirmationType `json:"confirmationType,omitempty"`
	// Only appears when ConfirmationType is "message"
	// RichText, can be Slate or Lexical
	ConfirmationMessage any `json:"confirmationMessage,omitempty"`
	// Only appears when ConfirmationType is "redirect"
	Redirect *FormRedirect `json:"redirect,omitempty"`
	// Optional
	Emails    []FormEmail `json:"emails,omitempty"`
	UpdatedAt string      `json:"updatedAt"`
	CreatedAt string      `json:"createdAt"`
}

// FormBlockType defines the type of field within a form.
type FormBlockType string

const (
	FormBlockTypeCheckbox FormBlockType = "checkbox"
	FormBlockTypeCountry  FormBlockType = "country"
	FormBlockTypeEmail    FormBlockType = "email"
	FormBlockTypeMessage  FormBlockType = "message"
	FormBlockTypeNumber   FormBlockType = "number"
	FormBlockTypeSelect   FormBlockType = "select"
	FormBlockTypeState    FormBlockType = "state"
	FormBlockTypeText     FormBlockType = "text"
	FormBlockTypeTextarea FormBlockType = "textarea"
	FormBlockTypeCurrency FormBlockType = "currency" // Not currently supported
)

// FormConfirmationType defines the way in which a user is directed after submitting a form.
type FormConfirmationType string

const (
	FormConfirmationTypeMessage  FormConfirmationType = "message"
	FormConfirmationTypeRedirect FormConfirmationType = "redirect"
)

// FormField represents a field in the Payload form builder.
type FormField struct {
	// Tabs that appear in all block types.
	ID        string        `json:"id,omitempty"`
	BlockType FormBlockType `json:"blockType"`
	Name      string        `json:"name"`
	Label     *string       `json:"label,omitempty"`
	Width     *int          `json:"width,omitempty"`
	Required  *bool         `json:"required,omitempty"`
	// One of the following fields must be present, depending on blockType
	DefaultValue *bool                    `json:"defaultValue,omitempty"`
	BlockName    *string                  `json:"blockName,omitempty"`
	Message      []map[string]interface{} `json:"message,omitempty"`
	Options      []FormOption             `json:"options,omitempty"`
}

// FormRedirect defines the type of confirmation message to display after
type FormRedirect struct {
	URL string `json:"url"`
}

// FormOption defines a singular option within a select field.
type FormOption struct {
	ID    *string `json:"id,omitempty"`
	Label string  `json:"label"`
	Value string  `json:"value"`
}

// FormEmail represents an email configuration for a form.
type FormEmail struct {
	ID        *string          `json:"id,omitempty"`
	EmailTo   *string          `json:"emailTo,omitempty"`
	CC        *string          `json:"cc,omitempty"`
	BCC       *string          `json:"bcc,omitempty"`
	ReplyTo   *string          `json:"replyTo,omitempty"`
	EmailFrom *string          `json:"emailFrom,omitempty"`
	Subject   string           `json:"subject"`
	Message   []map[string]any `json:"message,omitempty"`
}

// Render renders the form block to the provided writer as
// a form element.
func (f *Form) Render(_ context.Context, w io.Writer) error {
	return tpl.Templates.ExecuteTemplate(w, "form.html", f)
}
