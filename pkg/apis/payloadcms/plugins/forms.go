package payloadplugins

// Form Builder Plugin
// https://payloadcms.com/docs/plugins/form-builder
// TODO: Payment not currently supported

// Form defines a singular form collection type in the Form Builder Plugin
// within Payload CMS.
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

// FormSubmission defines a singular submission to a form in the Form Builder Plugin
// within Payload CMS.
type FormSubmission struct {
	ID             int                       `json:"id"`
	Form           interface{}               `json:"form"`
	SubmissionData []FormSubmissionDataEntry `json:"submissionData,omitempty"`
	UpdatedAt      string                    `json:"updatedAt"`
	CreatedAt      string                    `json:"createdAt"`
}

// FormSubmissionDataEntry defines a data entry within a form submission.
// These directly correspond to the fields within a form.
type FormSubmissionDataEntry struct {
	ID    *string `json:"id,omitempty"`
	Field string  `json:"field"`
	Value string  `json:"value"`
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
)

// FormConfirmationType defines the way in which a user is directed after submitting a form.
type FormConfirmationType string

const (
	FormConfirmationTypeMessage  FormConfirmationType = "message"
	FormConfirmationTypeRedirect FormConfirmationType = "redirect"
)

// FormField represents a field in the Payload form builder.
type FormField struct {
	// Fields that appear in all block types.
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

/**
 	High level example:
	{
		"id": 1,
		"title": "My Form",
		"fields": [
			{
				"id": "66544ca23d2ccf267279c9e6",
				"name": "name",
				"label": "Label",
				"width": 10,
				"required": null,
				"defaultValue": true,
				"blockName": null,
				"blockType": "checkbox"
			}
		],
		"submitButtonLabel": null,
		"confirmationType": "message",
		"confirmationMessage": [
			{
				"children": [
					{
						"text": "jkhkghjghjkkghj"
					}
				]
			}
		],
		"redirect": {
			"url": null
		},
		"emails": [],
		"updatedAt": "2024-05-27T09:04:34.702Z",
		"createdAt": "2024-05-27T09:04:34.702Z"
	}
*/
