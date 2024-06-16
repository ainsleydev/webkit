package payload

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestBlock_Decode(t *testing.T) {
	ttt := map[string]struct {
		input      Block
		decodeInto any
		want       any
		wantErr    bool
	}{
		"Decode into map[string]any": {
			input: Block{
				RawJSON: json.RawMessage(`{"field":"value"}`),
			},
			decodeInto: &map[string]any{},
			want:       &map[string]any{"field": "value"},
			wantErr:    false,
		},
		"Decode into struct": {
			input: Block{
				RawJSON: json.RawMessage(`{"field1":"value1","field2":2}`),
			},
			decodeInto: &struct {
				Field1 string
				Field2 int
			}{},
			want: &struct {
				Field1 string
				Field2 int
			}{Field1: "value1", Field2: 2},
			wantErr: false,
		},
		"Invalid JSON for struct": {
			input: Block{
				RawJSON: json.RawMessage(`{"field1":"value1","field2":"invalid_int"}`),
			},
			decodeInto: &struct {
				Field1 string
				Field2 int
			}{},
			want: &struct {
				Field1 string
				Field2 int
			}{},
			wantErr: true,
		},
		"Invalid JSON format": {
			input: Block{
				RawJSON: json.RawMessage(`{"field1":"value1","field2":2`),
			},
			decodeInto: &map[string]any{},
			want:       &map[string]any{},
			wantErr:    true,
		},
	}

	for name, test := range ttt {
		t.Run(name, func(t *testing.T) {
			err := test.input.Decode(test.decodeInto)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, test.decodeInto)
			}
		})
	}
}

func TestBlock_UnmarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input   string
		want    Block
		wantErr bool
	}{
		"All Fields Present": {
			input: `{"id":"123","blockType":"text","blockName":"Header","field":"value"}`,
			want: Block{
				Id:        "123",
				BlockType: "text",
				BlockName: ptr.StringPtr("Header"),
				Fields: map[string]any{
					"field": "value",
				},
				RawJSON: json.RawMessage(`{"id":"123","blockType":"text","blockName":"Header","field":"value"}`),
			},
			wantErr: false,
		},
		"Missing Optional Fields": {
			input: `{"id":"123","blockType":"text","field":"value"}`,
			want: Block{
				Id:        "123",
				BlockType: "text",
				BlockName: nil,
				Fields: map[string]any{
					"field": "value",
				},
				RawJSON: json.RawMessage(`{"id":"123","blockType":"text","field":"value"}`),
			},
			wantErr: false,
		},
		"Additional Fields": {
			input: `{"id":"123","blockType":"text","blockName":"Header","extra":"extraValue"}`,
			want: Block{
				Id:        "123",
				BlockType: "text",
				BlockName: ptr.StringPtr("Header"),
				Fields: map[string]any{
					"extra": "extraValue",
				},
				RawJSON: json.RawMessage(`{"id":"123","blockType":"text","blockName":"Header","extra":"extraValue"}`),
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `wrong`,
			want:    Block{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var b Block
			err := json.Unmarshal([]byte(test.input), &b)
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, b)
		})
	}
}
