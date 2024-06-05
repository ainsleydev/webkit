package payload

import "github.com/mitchellh/mapstructure"

// Block defines a common structure of a singular block layout
// from PayloadCMS.
type Block struct {
	BlockType any     `json:"block_type"`
	BlockName *string `json:"block_name,omitempty"`
	Id        *string `json:"id,omitempty"`
}

// https://stackoverflow.com/questions/33436730/unmarshal-json-with-some-known-and-some-unknown-field-names

// Blocks is a collection of Block types.
type Blocks []Block

// Decode decodes the input into the Blocks type.
func (b *Blocks) Decode(in any) error {
	return mapstructure.Decode(in, &b)
}
