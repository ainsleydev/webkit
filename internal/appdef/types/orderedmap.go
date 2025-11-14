package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// OrderedMap preserves insertion order when marshaling/unmarshaling JSON.
// This ensures that JSON object keys maintain their original order rather than
// being alphabetically sorted during marshal operations.
type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

// NewOrderedMap creates a new OrderedMap.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:   []K{},
		values: make(map[K]V),
	}
}

// Get retrieves a value by key.
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	if om == nil || om.values == nil {
		var zero V
		return zero, false
	}
	v, ok := om.values[key]
	return v, ok
}

// Set stores a value with the given key, preserving insertion order.
func (om *OrderedMap[K, V]) Set(key K, value V) {
	if om.values == nil {
		om.values = make(map[K]V)
	}
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

// Keys returns all keys in insertion order.
func (om *OrderedMap[K, V]) Keys() []K {
	if om == nil {
		return nil
	}
	return om.keys
}

// Len returns the number of key-value pairs.
func (om *OrderedMap[K, V]) Len() int {
	if om == nil || om.values == nil {
		return 0
	}
	return len(om.values)
}

// Range iterates over the map in insertion order.
func (om *OrderedMap[K, V]) Range(fn func(key K, value V) bool) {
	if om == nil {
		return
	}
	for _, key := range om.keys {
		if !fn(key, om.values[key]) {
			break
		}
	}
}

// UnmarshalJSON implements json.Unmarshaler, preserving key order.
func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	if om.values == nil {
		om.values = make(map[K]V)
	}
	om.keys = []K{}

	// First unmarshal into a temporary map to get key-value pairs
	var temp map[K]V
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Parse the JSON again to extract key order
	// We'll use a decoder and manually parse to maintain order
	dec := json.NewDecoder(bytes.NewReader(data))

	// Read opening brace
	token, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected object start, got %v", token)
	}

	// Read key-value pairs in order
	for dec.More() {
		// Read key
		token, err := dec.Token()
		if err != nil {
			return err
		}

		// Convert token to key type
		var key K
		keyBytes, err := json.Marshal(token)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(keyBytes, &key); err != nil {
			return fmt.Errorf("failed to convert key: %w", err)
		}

		// Read value (skip it since we already have it in temp map)
		var rawValue json.RawMessage
		if err := dec.Decode(&rawValue); err != nil {
			return err
		}

		// Store in order
		om.keys = append(om.keys, key)
	}

	// Read closing brace
	token, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected object end, got %v", token)
	}

	// Now copy values from temp map in the order we discovered
	for _, key := range om.keys {
		om.values[key] = temp[key]
	}

	return nil
}

// MarshalJSON implements json.Marshaler, writing keys in insertion order.
func (om *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	if om.values == nil || len(om.keys) == 0 {
		return []byte("{}"), nil
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')

	first := true
	for _, key := range om.keys {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		// Marshal key
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		// Marshal value
		valueBytes, err := json.Marshal(om.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valueBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}
