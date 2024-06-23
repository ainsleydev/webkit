package payload

import (
	"testing"

	"github.com/goccy/go-json"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestBlock_Decode(t *testing.T) {
	type decodeInto struct {
		Field1 string
		Field2 int
	}

	tt := map[string]struct {
		input   Block
		want    any
		wantErr bool
	}{
		"Decode OK": {
			input: Block{
				RawJSON: json.RawMessage(`{"field1":"value1","field2":2}`),
			},
			want: decodeInto{
				Field1: "value1",
				Field2: 2,
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input: Block{
				RawJSON: json.RawMessage(`{wrong}`),
			},
			want:    decodeInto{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var d decodeInto
			err := test.input.Decode(&d)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, d)
		})
	}
}

func TestBlock_UnmarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input   string
		want    Block
		wantErr bool
	}{
		"OK": {
			input: `{"id":"123","blockType":"text","blockName":"Header","field":"value"}`,
			want: Block{
				ID:        "123",
				BlockType: "text",
				BlockName: ptr.StringPtr("Header"),
				Fields: map[string]any{
					"field": "value",
				},
				RawJSON: json.RawMessage(`{"id":"123","blockType":"text","blockName":"Header","field":"value"}`),
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{block_name:2345}`,
			want:    Block{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var b Block
			err := b.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, b)
		})
	}
}
