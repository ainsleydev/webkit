package payload

import (
	"errors"

	"github.com/goccy/go-json"
)

// JSON represents a map that can be marshaled  into a
// Payload JSON field.
//
// Payload expects it to be created as a string.
type JSON map[string]any

// MarshalJSON marshals the JSON map to a string.
//
//goland:noinspection GoMixedReceiverTypes
func (j JSON) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(map[string]any(j))
	if err != nil {
		return nil, errors.New("payload.JSON: " + err.Error())
	}
	return data, nil
}

// UnmarshalJSON unmarshals a string into a JSON map.
//
//goland:noinspection GoMixedReceiverTypes
func (j *JSON) UnmarshalJSON(data []byte) error {
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return errors.New("payload.JSON: " + err.Error())

	}
	*j = obj
	return nil
}
