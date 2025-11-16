package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/jsonschema-go"
)

func TestOrderedMap_PreservesInsertionOrder(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()
	om.Set("zebra", 1)
	om.Set("apple", 2)
	om.Set("mango", 3)
	om.Set("banana", 4)

	keys := om.Keys()
	assert.Equal(t, []string{"zebra", "apple", "mango", "banana"}, keys)
}

func TestOrderedMap_MarshalJSON_PreservesOrder(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, string]()
	om.Set("generate", "go generate ./...")
	om.Set("build", "go build main.go")
	om.Set("format", "go fmt ./...")
	om.Set("lint", "echo")
	om.Set("test", "go test ./...")

	data, err := json.Marshal(om)
	require.NoError(t, err)

	expected := `{"generate":"go generate ./...","build":"go build main.go","format":"go fmt ./...","lint":"echo","test":"go test ./..."}`
	assert.JSONEq(t, expected, string(data))

	// Verify exact order (not just JSON equality)
	assert.Equal(t, expected, string(data))
}

func TestOrderedMap_UnmarshalJSON_PreservesOrder(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{"generate":"go generate ./...","build":"go build main.go","format":"go fmt ./...","lint":"echo","test":"go test ./..."}`)

	om := NewOrderedMap[string, string]()
	err := json.Unmarshal(jsonData, om)
	require.NoError(t, err)

	keys := om.Keys()
	assert.Equal(t, []string{"generate", "build", "format", "lint", "test"}, keys)

	val, ok := om.Get("generate")
	assert.True(t, ok)
	assert.Equal(t, "go generate ./...", val)

	val, ok = om.Get("build")
	assert.True(t, ok)
	assert.Equal(t, "go build main.go", val)
}

func TestOrderedMap_RoundTrip_PreservesOrder(t *testing.T) {
	t.Parallel()

	// Original JSON with specific order
	original := []byte(`{"first":"1","second":"2","third":"3","fourth":"4"}`)

	// Unmarshal
	om := NewOrderedMap[string, string]()
	err := json.Unmarshal(original, om)
	require.NoError(t, err)

	// Marshal back
	result, err := json.Marshal(om)
	require.NoError(t, err)

	// Should maintain exact order
	assert.Equal(t, string(original), string(result))
}

func TestOrderedMap_Get(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()
	om.Set("key1", 100)
	om.Set("key2", 200)

	val, ok := om.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	val, ok = om.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, 200, val)

	val, ok = om.Get("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

func TestOrderedMap_Set_UpdateExisting(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()
	om.Set("key1", 100)
	om.Set("key2", 200)
	om.Set("key1", 300) // Update existing

	// Should maintain original order
	keys := om.Keys()
	assert.Equal(t, []string{"key1", "key2"}, keys)

	// Should have updated value
	val, ok := om.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 300, val)
}

func TestOrderedMap_Range(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()
	om.Set("first", 1)
	om.Set("second", 2)
	om.Set("third", 3)

	var keys []string
	var values []int
	om.Range(func(key string, value int) bool {
		keys = append(keys, key)
		values = append(values, value)
		return true
	})

	assert.Equal(t, []string{"first", "second", "third"}, keys)
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestOrderedMap_Range_EarlyExit(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()
	om.Set("first", 1)
	om.Set("second", 2)
	om.Set("third", 3)

	var keys []string
	om.Range(func(key string, value int) bool {
		keys = append(keys, key)
		return key != "second" // Stop after second
	})

	assert.Equal(t, []string{"first", "second"}, keys)
}

func TestOrderedMap_Len(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()
	assert.Equal(t, 0, om.Len())

	om.Set("key1", 1)
	assert.Equal(t, 1, om.Len())

	om.Set("key2", 2)
	assert.Equal(t, 2, om.Len())

	om.Set("key1", 100) // Update shouldn't change length
	assert.Equal(t, 2, om.Len())
}

func TestOrderedMap_Empty(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()
	data, err := json.Marshal(om)
	require.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestOrderedMap_Nil(t *testing.T) {
	t.Parallel()

	var om *OrderedMap[string, int]
	assert.Equal(t, 0, om.Len())
	assert.Nil(t, om.Keys())

	val, ok := om.Get("key")
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	// Range should not panic on nil
	om.Range(func(key string, value int) bool {
		t.Fatal("should not be called on nil map")
		return true
	})
}

func TestOrderedMap_Set_NilValues(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		setup func() *OrderedMap[string, int]
		key   string
		value int
		want  []string
	}{
		"Set on uninitialized values map": {
			setup: func() *OrderedMap[string, int] {
				return &OrderedMap[string, int]{
					keys:   []string{},
					values: nil, // Explicitly nil
				}
			},
			key:   "first",
			value: 100,
			want:  []string{"first"},
		},
		"Set multiple on nil values": {
			setup: func() *OrderedMap[string, int] {
				return &OrderedMap[string, int]{
					keys:   []string{},
					values: nil,
				}
			},
			key:   "second",
			value: 200,
			want:  []string{"second"},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			om := test.setup()
			om.Set(test.key, test.value)

			assert.Equal(t, test.want, om.Keys())
			val, ok := om.Get(test.key)
			assert.True(t, ok)
			assert.Equal(t, test.value, val)
		})
	}
}

func TestOrderedMap_UnmarshalJSON_Errors(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		wantErr bool
	}{
		"Invalid JSON": {
			input:   `{invalid json}`,
			wantErr: true,
		},
		"Not an object": {
			input:   `["array", "not", "object"]`,
			wantErr: true,
		},
		"Unclosed object": {
			input:   `{"key": "value"`,
			wantErr: true,
		},
		"Invalid key type for string keys": {
			input:   `{123: "value"}`,
			wantErr: true,
		},
		"Missing value": {
			input:   `{"key":}`,
			wantErr: true,
		},
		"Invalid value JSON": {
			input:   `{"key": invalid}`,
			wantErr: true,
		},
		"Extra data after object": {
			input:   `{"key": "value"} extra`,
			wantErr: true,
		},
		"Nested object with error": {
			input:   `{"key": {"nested": invalid}}`,
			wantErr: true,
		},
		"Valid empty object": {
			input:   `{}`,
			wantErr: false,
		},
		"Valid single entry": {
			input:   `{"key": "value"}`,
			wantErr: false,
		},
		"Valid multiple entries": {
			input:   `{"a": "1", "b": "2", "c": "3"}`,
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			om := NewOrderedMap[string, string]()
			err := json.Unmarshal([]byte(test.input), om)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderedMap_UnmarshalJSON_InvalidValueDecode(t *testing.T) {
	t.Parallel()

	// Test case where the value cannot be decoded
	type CustomType struct {
		Value string `json:"value"`
	}

	jsonData := []byte(`{"key1": {"value": "valid"}, "key2": "invalid_type"}`)
	om := NewOrderedMap[string, CustomType]()
	err := json.Unmarshal(jsonData, om)

	// This should error because "invalid_type" string can't unmarshal into CustomType
	assert.Error(t, err)
}

func TestOrderedMap_MarshalJSON_NilValues(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		om   *OrderedMap[string, string]
		want string
	}{
		"Nil values map": {
			om: &OrderedMap[string, string]{
				keys:   []string{},
				values: nil,
			},
			want: "{}",
		},
		"Empty keys": {
			om: &OrderedMap[string, string]{
				keys:   []string{},
				values: make(map[string]string),
			},
			want: "{}",
		},
		"Nil keys with values": {
			om: &OrderedMap[string, string]{
				keys:   nil,
				values: map[string]string{"key": "value"},
			},
			want: "{}",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(test.om)
			require.NoError(t, err)
			assert.Equal(t, test.want, string(data))
		})
	}
}

func TestOrderedMap_ComplexTypes(t *testing.T) {
	t.Parallel()

	t.Run("Struct values", func(t *testing.T) {
		t.Parallel()

		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		om := NewOrderedMap[string, TestStruct]()
		om.Set("first", TestStruct{Name: "First", Value: 1})
		om.Set("second", TestStruct{Name: "Second", Value: 2})

		keys := om.Keys()
		assert.Equal(t, []string{"first", "second"}, keys)

		data, err := json.Marshal(om)
		require.NoError(t, err)
		assert.JSONEq(t, `{"first":{"name":"First","value":1},"second":{"name":"Second","value":2}}`, string(data))
	})
}

func TestOrderedMap_LargeDataSet(t *testing.T) {
	t.Parallel()

	om := NewOrderedMap[string, int]()

	// Insert 100 items
	for i := 0; i < 100; i++ {
		om.Set(fmt.Sprintf("key-%d", i), i)
	}

	assert.Equal(t, 100, om.Len())

	// Verify order is preserved
	keys := om.Keys()
	for i := 0; i < 100; i++ {
		assert.Equal(t, fmt.Sprintf("key-%d", i), keys[i])
		val, ok := om.Get(keys[i])
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}

	// Marshal and unmarshal
	data, err := json.Marshal(om)
	require.NoError(t, err)

	om2 := NewOrderedMap[string, int]()
	err = json.Unmarshal(data, om2)
	require.NoError(t, err)

	// Verify order is still preserved
	keys2 := om2.Keys()
	assert.Equal(t, keys, keys2)
}

func TestOrderedMap_NestedValues(t *testing.T) {
	t.Parallel()

	type NestedStruct struct {
		Items map[string]int `json:"items"`
		Count int            `json:"count"`
	}

	om := NewOrderedMap[string, NestedStruct]()
	om.Set("first", NestedStruct{
		Items: map[string]int{"a": 1, "b": 2},
		Count: 2,
	})
	om.Set("second", NestedStruct{
		Items: map[string]int{"c": 3},
		Count: 1,
	})

	// Verify order
	keys := om.Keys()
	assert.Equal(t, []string{"first", "second"}, keys)

	// Marshal and verify structure
	data, err := json.Marshal(om)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"first"`)
	assert.Contains(t, string(data), `"second"`)

	// Unmarshal and verify order is preserved
	om2 := NewOrderedMap[string, NestedStruct]()
	err = json.Unmarshal(data, om2)
	require.NoError(t, err)
	assert.Equal(t, keys, om2.Keys())
}

func TestOrderedMap_MarshalJSON_Errors(t *testing.T) {
	t.Parallel()

	t.Run("Unmarshalable value type", func(t *testing.T) {
		t.Parallel()

		// Create a type that fails to marshal
		type FailingMarshaler struct{}

		om := NewOrderedMap[string, FailingMarshaler]()
		om.Set("key", FailingMarshaler{})

		// Marshal should fail because FailingMarshaler contains unexported fields
		// that json package cannot handle (channels, funcs, etc)
		type BadValue struct {
			Ch chan int // channels cannot be marshaled
		}

		om2 := NewOrderedMap[string, BadValue]()
		om2.Set("key", BadValue{Ch: make(chan int)})

		_, err := json.Marshal(om2)
		assert.Error(t, err)
	})
}

func TestOrderedMap_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("Set and overwrite preserves position", func(t *testing.T) {
		t.Parallel()

		om := NewOrderedMap[string, string]()
		om.Set("a", "1")
		om.Set("b", "2")
		om.Set("c", "3")
		om.Set("b", "updated")

		keys := om.Keys()
		assert.Equal(t, []string{"a", "b", "c"}, keys)
		val, _ := om.Get("b")
		assert.Equal(t, "updated", val)
	})

	t.Run("Empty string keys and values", func(t *testing.T) {
		t.Parallel()

		om := NewOrderedMap[string, string]()
		om.Set("", "")
		om.Set("a", "")
		om.Set("", "value")

		assert.Equal(t, 2, om.Len())
		val, ok := om.Get("")
		assert.True(t, ok)
		assert.Equal(t, "value", val)
	})

	t.Run("Marshal and unmarshal empty strings", func(t *testing.T) {
		t.Parallel()

		om := NewOrderedMap[string, string]()
		om.Set("", "empty key")
		om.Set("normal", "")

		data, err := json.Marshal(om)
		require.NoError(t, err)

		om2 := NewOrderedMap[string, string]()
		err = json.Unmarshal(data, om2)
		require.NoError(t, err)

		assert.Equal(t, om.Keys(), om2.Keys())
	})
}

func TestOrderedMap_JSONSchema(t *testing.T) {
	t.Parallel()

	t.Run("String values", func(t *testing.T) {
		t.Parallel()

		om := NewOrderedMap[string, string]()
		schema, err := om.JSONSchema()
		require.NoError(t, err)

		t.Log("Verify schema type is object")
		{
			assert.NotNil(t, schema.Type)
			assert.NotNil(t, schema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *schema.Type.SimpleTypes)
		}

		t.Log("Verify additionalProperties is set")
		{
			assert.NotNil(t, schema.AdditionalProperties)
			assert.NotNil(t, schema.AdditionalProperties.TypeObject)
		}

		t.Log("Verify additionalProperties schema for string type")
		{
			valueSchema := schema.AdditionalProperties.TypeObject
			assert.NotNil(t, valueSchema.Type)
			assert.NotNil(t, valueSchema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.String, *valueSchema.Type.SimpleTypes)
		}
	})

	t.Run("Integer values", func(t *testing.T) {
		t.Parallel()

		om := NewOrderedMap[string, int]()
		schema, err := om.JSONSchema()
		require.NoError(t, err)

		t.Log("Verify schema type is object")
		{
			assert.NotNil(t, schema.Type)
			assert.NotNil(t, schema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *schema.Type.SimpleTypes)
		}

		t.Log("Verify additionalProperties schema for int type")
		{
			assert.NotNil(t, schema.AdditionalProperties)
			assert.NotNil(t, schema.AdditionalProperties.TypeObject)
			valueSchema := schema.AdditionalProperties.TypeObject
			assert.NotNil(t, valueSchema.Type)
			assert.NotNil(t, valueSchema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Integer, *valueSchema.Type.SimpleTypes)
		}
	})

	t.Run("Complex struct values", func(t *testing.T) {
		t.Parallel()

		type TestStruct struct {
			Command string `json:"command"`
			SkipCI  bool   `json:"skip_ci"`
			Timeout string `json:"timeout"`
		}

		om := NewOrderedMap[string, TestStruct]()
		schema, err := om.JSONSchema()
		require.NoError(t, err)

		t.Log("Verify schema type is object")
		{
			assert.NotNil(t, schema.Type)
			assert.NotNil(t, schema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *schema.Type.SimpleTypes)
		}

		t.Log("Verify additionalProperties contains struct schema")
		{
			assert.NotNil(t, schema.AdditionalProperties)
			assert.NotNil(t, schema.AdditionalProperties.TypeObject)

			valueSchema := schema.AdditionalProperties.TypeObject
			assert.NotNil(t, valueSchema.Type)
			assert.NotNil(t, valueSchema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *valueSchema.Type.SimpleTypes)
		}

		t.Log("Verify struct properties are in schema")
		{
			valueSchema := schema.AdditionalProperties.TypeObject
			assert.NotNil(t, valueSchema.Properties)
			assert.Contains(t, valueSchema.Properties, "command")
			assert.Contains(t, valueSchema.Properties, "skip_ci")
			assert.Contains(t, valueSchema.Properties, "timeout")
		}
	})

	t.Run("Boolean values", func(t *testing.T) {
		t.Parallel()

		om := NewOrderedMap[string, bool]()
		schema, err := om.JSONSchema()
		require.NoError(t, err)

		t.Log("Verify schema type is object")
		{
			assert.NotNil(t, schema.Type)
			assert.NotNil(t, schema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *schema.Type.SimpleTypes)
		}

		t.Log("Verify additionalProperties schema for bool type")
		{
			assert.NotNil(t, schema.AdditionalProperties)
			valueSchema := schema.AdditionalProperties.TypeObject
			assert.NotNil(t, valueSchema)
			assert.NotNil(t, valueSchema.Type)
			assert.NotNil(t, valueSchema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Boolean, *valueSchema.Type.SimpleTypes)
		}
	})

	t.Run("Nil OrderedMap", func(t *testing.T) {
		t.Parallel()

		var om *OrderedMap[string, string]
		schema, err := om.JSONSchema()
		require.NoError(t, err)

		t.Log("Should still return valid object schema")
		{
			assert.NotNil(t, schema.Type)
			assert.NotNil(t, schema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *schema.Type.SimpleTypes)
		}
	})

	t.Run("Empty OrderedMap", func(t *testing.T) {
		t.Parallel()

		om := NewOrderedMap[string, string]()
		schema, err := om.JSONSchema()
		require.NoError(t, err)

		t.Log("Empty map should still generate proper schema")
		{
			assert.NotNil(t, schema.Type)
			assert.NotNil(t, schema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *schema.Type.SimpleTypes)
			assert.NotNil(t, schema.AdditionalProperties)
		}
	})

	t.Run("Nested struct with pointer", func(t *testing.T) {
		t.Parallel()

		type NestedConfig struct {
			Name  string `json:"name"`
			Value *int   `json:"value,omitempty"`
		}

		om := NewOrderedMap[string, NestedConfig]()
		schema, err := om.JSONSchema()
		require.NoError(t, err)

		t.Log("Verify schema structure for nested struct")
		{
			assert.NotNil(t, schema.Type)
			assert.NotNil(t, schema.Type.SimpleTypes)
			assert.Equal(t, jsonschema.Object, *schema.Type.SimpleTypes)
			assert.NotNil(t, schema.AdditionalProperties)
			assert.NotNil(t, schema.AdditionalProperties.TypeObject)
		}
	})
}
