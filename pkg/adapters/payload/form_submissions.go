package payload

import "time"

// FormSubmission defines a singular submission to a form in the Form Builder Plugin
// within Payload CMS.
type FormSubmission struct {
	ID             int                       `json:"id"`
	Form           Form                      `json:"form"`
	SubmissionData []FormSubmissionDataEntry `json:"submissionData,omitempty"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
	CreatedAt      time.Time                 `json:"createdAt"`
}

// FormSubmissionDataEntry defines a data entry within a form submission.
// These directly correspond to the fields within a form.
type FormSubmissionDataEntry struct {
	ID    *string `json:"id,omitempty"`
	Field string  `json:"field"`
	Value string  `json:"value"`
}
