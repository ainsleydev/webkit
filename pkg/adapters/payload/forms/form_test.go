package payloadforms

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

var dummyForm = Form{
	ID:    1,
	Title: "Contact Form",
	Fields: []FormField{
		{
			ID:        "f1",
			Name:      "name",
			Label:     ptr.StringPtr("Your Name"),
			BlockType: FormBlockTypeText,
			Required:  ptr.BoolPtr(true),
		},
		{
			ID:        "f2",
			Name:      "email",
			Label:     ptr.StringPtr("Email Address"),
			BlockType: FormBlockTypeEmail,
			Required:  ptr.BoolPtr(true),
		},
		{
			ID:        "f3",
			Name:      "message",
			Label:     ptr.StringPtr("Message"),
			BlockType: FormBlockTypeTextarea,
			Required:  ptr.BoolPtr(true),
		},
	},
	// SubmitButtonLabel: stringPtr("Submit"),
	ConfirmationType: FormConfirmationTypeMessage,
}

func TestFormFields(t *testing.T) {
	buf := &bytes.Buffer{}
	PayloadFormFields(dummyForm.Fields).Render(context.TODO(), buf)
	fmt.Println(buf.String())
}
