package cache

import (
	"encoding/json"
)

// marshalIdent marshals a JSON byte slice with indentation used
// for pretty printing in OS and File cache stores.
func marshalIdent(b []byte) ([]byte, error) {
	var jsonData any
	if err := json.Unmarshal(b, &jsonData); err != nil {
		return nil, err
	}
	return json.MarshalIndent(jsonData, "", "\t")
}

// contains checks if a string is in a slice of strings.
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
