package main

type FieldValues map[string]interface{}

type PaymentFieldConfig struct {
	Field            *Field
	PaymentProcessor *SelectField
}

type FieldConfig interface{}

type TextField struct {
	BlockName    *string
	BlockType    string
	DefaultValue *string
	Label        *string
	Name         string
	Required     *bool
	Width        *int
}

type TextAreaField struct {
	BlockName    *string
	BlockType    string
	DefaultValue *string
	Label        *string
	Name         string
	Required     *bool
	Width        *int
}

type SelectFieldOption struct {
	Label string
	Value string
}

type SelectField struct {
	BlockName    *string
	BlockType    string
	DefaultValue *string
	Label        *string
	Name         string
	Options      []SelectFieldOption
	Required     *bool
	Width        *int
}

type PriceCondition struct {
	Condition         string
	FieldToUse        string
	Operator          string
	ValueForCondition string
	ValueForOperator  interface{}
	ValueType         string
}

type PaymentField struct {
	BasePrice        float64
	BlockName        *string
	BlockType        string
	DefaultValue     *string
	Label            *string
	Name             string
	PaymentProcessor string
	PriceConditions  []PriceCondition
	Required         *bool
	Width            *int
}

type EmailField struct {
	BlockName    *string
	BlockType    string
	DefaultValue *string
	Label        *string
	Name         string
	Required     *bool
	Width        *int
}

type StateField struct {
	BlockName    *string
	BlockType    string
	DefaultValue *string
	Label        *string
	Name         string
	Required     *bool
	Width        *int
}

type CountryField struct {
	BlockName    *string
	BlockType    string
	DefaultValue *string
	Label        *string
	Name         string
	Required     *bool
	Width        *int
}

type CheckboxField struct {
	BlockName    *string
	BlockType    string
	DefaultValue *bool
	Label        *string
	Name         string
	Required     *bool
	Width        *int
}

type MessageField struct {
	BlockName *string
	BlockType string
	Message   interface{}
}

type FormFieldBlock interface {
	TextField
	TextAreaField
	SelectField
}

type Email struct {
	Bcc       *string
	Cc        *string
	EmailFrom string
	EmailTo   string
	Message   interface{}
	ReplyTo   *string
	Subject   string
}

type Redirect struct {
	// Add any necessary fields for the redirect
}

type Form struct {
	ConfirmationMessage *interface{}
	ConfirmationType    string
	Emails              []Email
	Fields              []FormFieldBlock
	ID                  string
	Redirect            *Redirect
	SubmitButtonLabel   *string
	Title               string
}
