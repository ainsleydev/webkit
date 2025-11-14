package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
