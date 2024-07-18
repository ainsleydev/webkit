package payload

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

var form = ` {
	"id": 1,
	"content": "Test",
	"title": "Contact",
	"fields": [
		{
			"id": "66978fa9cfde9d48af3ee2ea",
			"name": "name",
			"label": "Name",
			"width": null,
			"defaultValue": null,
			"required": true,
			"blockName": null,
			"blockType": "text"
		},
		{
			"id": "66978fa9cfde9d48af3ee2eb",
			"name": "email",
			"label": "Email",
			"width": null,
			"defaultValue": null,
			"required": true,
			"blockName": null,
			"blockType": "text"
		},
		{
			"id": "66978fa9cfde9d48af3ee2ec",
			"name": "message",
			"label": "Message",
			"width": null,
			"defaultValue": null,
			"required": null,
			"blockName": null,
			"blockType": "textarea"
		}
	],
	"submitButtonLabel": "Get In Touch",
	"redirect": {
		"url": null
	},
	"emails": null,
	"confirmationMessage": null,
	"updatedAt": "2024-07-17T09:32:25.293Z",
	"createdAt": "2024-07-17T09:32:25.293Z"
}`

func TestForm_UnmarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input   string
		want    Form
		wantErr bool
	}{
		"OK": {
			input: form,
			want: Form{
				ID:    1,
				Title: "Contact",
				Fields: []FormField{
					{
						ID:        "66978fa9cfde9d48af3ee2ea",
						Name:      "name",
						Label:     "Name",
						Required:  ptr.BoolPtr(true),
						BlockType: FormBlockTypeText,
					},
					{
						ID:        "66978fa9cfde9d48af3ee2eb",
						Name:      "email",
						Label:     "Email",
						Required:  ptr.BoolPtr(true),
						BlockType: FormBlockTypeText,
					},
					{
						ID:        "66978fa9cfde9d48af3ee2ec",
						Name:      "message",
						Label:     "Message",
						BlockType: FormBlockTypeTextarea,
					},
				},
				SubmitButtonLabel: ptr.StringPtr("Get In Touch"),
				Redirect:          &FormRedirect{URL: ""},
				Extra: map[string]any{
					"content": "Test",
				},
				UpdatedAt: "2024-07-17T09:32:25.293Z",
				CreatedAt: "2024-07-17T09:32:25.293Z",
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{wrong}`,
			want:    Form{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var f Form
			err := f.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, f)
		})
	}
}

func TestFormField_Render(t *testing.T) {
	ttt := map[string]struct {
		input FormField
		want  string
		err   error
	}{
		"Text Field": {
			input: FormField{
				ID:        "text-id",
				BlockType: FormBlockTypeText,
				Name:      "text-name",
			},
			want: `<input class="form-input" type="text" name="text-name" id="text-id" />`,
			err:  nil,
		},
		"Email Field": {
			input: FormField{
				ID:        "email-id",
				BlockType: FormBlockTypeEmail,
				Name:      "email-name",
			},
			want: `<input class="form-input" type="email" name="email-name" id="email-id" />`,
			err:  nil,
		},
		"Textarea Field": {
			input: FormField{
				ID:        "textarea-id",
				BlockType: FormBlockTypeTextarea,
				Name:      "textarea-name",
			},
			want: `<textarea class="form-input form-textarea" rows="6" name="textarea-name" id="textarea-id"></textarea>`,
			err:  nil,
		},
		"Not Found": {
			input: FormField{
				ID:        "unknown-id",
				BlockType: FormBlockType("unknown"),
				Name:      "unknown-name",
			},
			want: "",
			err:  fmt.Errorf("no renderer found for block type unknown"),
		},
	}

	for name, test := range ttt {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			err := test.input.Render(context.TODO(), &buf)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestFormBlockType_String(t *testing.T) {
	tt := map[string]struct {
		input *FormBlockType
		want  string
	}{
		"Nil FormBlockType": {
			input: nil,
			want:  "",
		},
		"Empty FormBlockType": {
			input: func() *FormBlockType {
				f := FormBlockType("")
				return &f
			}(),
			want: "",
		},
		"Non-empty FormBlockType": {
			input: func() *FormBlockType {
				f := FormBlockType("text")
				return &f
			}(),
			want: "text",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, test.input.String())
		})
	}
}

func TestFormBlockType_UnmarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input []byte
		want  FormBlockType
	}{
		"Empty JSON String": {
			input: []byte(`""`),
			want:  FormBlockType(""),
		},
		"Non-empty JSON String": {
			input: []byte(`"text"`),
			want:  FormBlockType("text"),
		},
		"JSON String with Spaces": {
			input: []byte(`"text with spaces"`),
			want:  FormBlockType("text with spaces"),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var got FormBlockType
			err := json.Unmarshal(test.input, &got)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}
