package bellows

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	B
	F     int
	Inner Inner
}

type Inner struct {
	B
	V string
}

type B struct {
	C string
	D int
}

func TestAnonStructFlatten(t *testing.T) {
	expecting := map[string]interface{}{
		"F":       1,
		"C":       "test",
		"D":       2,
		"Inner.C": "",
		"Inner.D": 0,
		"Inner.V": "",
	}

	a := A{
		F: 1,
		B: B{
			C: "test",
			D: 2,
		},
	}

	result := Flatten(a)
	assert.Equal(t, expecting, result)
}

func TestFlattenStructWithPointers(t *testing.T) {
	type User struct {
		Name *string
		Age  *int
	}
	name := "John"
	age := 30
	input := User{
		Name: &name,
		Age:  &age,
	}
	expected := map[string]interface{}{
		"Name": &name,
		"Age":  &age,
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenStructWithNilPointers(t *testing.T) {
	type User struct {
		Name *string
		Age  *int
	}
	input := User{
		Name: nil,
		Age:  nil,
	}
	expected := map[string]interface{}{
		"Name": nil,
		"Age":  nil,
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenStructWithUnexportedFields(t *testing.T) {
	type TestStruct struct {
		PublicField    string
		privateField   string
		AnotherPublic  int
		anotherPrivate bool
		NestedStruct   struct {
			PublicNested  string
			privateNested int
		}
	}

	input := TestStruct{
		PublicField:    "public_value",
		privateField:   "private_value",
		AnotherPublic:  42,
		anotherPrivate: true,
		NestedStruct: struct {
			PublicNested  string
			privateNested int
		}{
			PublicNested:  "nested_public",
			privateNested: 100,
		},
	}

	expected := map[string]interface{}{
		"PublicField":             "public_value",
		"AnotherPublic":           42,
		"NestedStruct.PublicNested": "nested_public",
	}

	// Test that flattening works fine with only exported fields
	assert.NotPanics(t, func() {
		result := Flatten(input)
		assert.Equal(t, expected, result)

		// Verify that all keys correspond to exported fields
		for key := range result {
			// All keys should start with uppercase (exported fields only)
			assert.True(t, key[0] >= 'A' && key[0] <= 'Z', "Found unexported field in result: %s", key)
		}
	})
}

func TestFlattenStructWithOnlyExportedFields(t *testing.T) {
	type SafeStruct struct {
		PublicField   string
		AnotherPublic int
		NestedStruct  struct {
			PublicNested string
		}
	}

	input := SafeStruct{
		PublicField:   "public_value",
		AnotherPublic: 42,
		NestedStruct: struct {
			PublicNested string
		}{
			PublicNested: "nested_public",
		},
	}

	expected := map[string]interface{}{
		"PublicField":             "public_value",
		"AnotherPublic":           42,
		"NestedStruct.PublicNested": "nested_public",
	}

	// Test that flattening works fine with only exported fields
	assert.NotPanics(t, func() {
		result := Flatten(input)
		assert.Equal(t, expected, result)

		// Verify that all keys correspond to exported fields
		for key := range result {
			// All keys should start with uppercase (exported fields only)
			assert.True(t, key[0] >= 'A' && key[0] <= 'Z', "Found unexported field in result: %s", key)
		}
	})
}

func TestFlattenStructWithMixedFieldTypesSafe(t *testing.T) {
	// Test with only exported fields to avoid panic
	type SafeComplexStruct struct {
		ExportedString string
		ExportedInt    int
		ExportedBool   bool
		ExportedSlice  []string
		ExportedMap    map[string]interface{}
	}

	input := SafeComplexStruct{
		ExportedString: "exported",
		ExportedInt:    100,
		ExportedBool:   true,
		ExportedSlice:  []string{"a", "b"},
		ExportedMap: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	expected := map[string]interface{}{
		"ExportedString":    "exported",
		"ExportedInt":       100,
		"ExportedBool":      true,
		"ExportedSlice.[0]": "a",
		"ExportedSlice.[1]": "b",
		"ExportedMap.key1":  "value1",
		"ExportedMap.key2":  42,
	}

	result := Flatten(input)
	assert.Equal(t, expected, result)

	// Verify all fields are properly flattened
	assert.Len(t, result, 7)
	assert.Contains(t, result, "ExportedString")
	assert.Contains(t, result, "ExportedInt")
	assert.Contains(t, result, "ExportedBool")
	assert.Contains(t, result, "ExportedSlice.[0]")
	assert.Contains(t, result, "ExportedSlice.[1]")
	assert.Contains(t, result, "ExportedMap.key1")
	assert.Contains(t, result, "ExportedMap.key2")
}