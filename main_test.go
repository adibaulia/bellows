package bellows

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	example = map[string]interface{}{
		"a": "b",
		"b": []string{"1", "2", "3"},
		"c": []interface{}{
			map[string]interface{}{"d": 1, "e": true, "k": []int{5, 6, 7}},
			map[string]interface{}{"d": 2, "e": false, "t": "112"},
		},
	}
	flat = Flatten(example)
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

func TestAll(t *testing.T) {
	assert := assert.New(t)

	flattened := Flatten(example)
	expanded := Expand(flattened)

	expecting, _ := json.Marshal(example)
	actual, _ := json.Marshal(expanded)
	assert.Equal(expecting, actual)
}

func TestFlattenBasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]interface{}
	}{
		{
			name:     "string value",
			input:    "hello",
			expected: map[string]interface{}{},
		},
		{
			name:     "integer value",
			input:    42,
			expected: map[string]interface{}{},
		},
		{
			name:     "boolean value",
			input:    true,
			expected: map[string]interface{}{},
		},
		{
			name:     "nil value",
			input:    nil,
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlattenSimpleMap(t *testing.T) {
	input := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
	}
	expected := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenNestedMap(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"details": map[string]interface{}{
				"age":  30,
				"city": "NYC",
			},
		},
	}
	expected := map[string]interface{}{
		"user.name":         "John",
		"user.details.age":  30,
		"user.details.city": "NYC",
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenArray(t *testing.T) {
	input := map[string]interface{}{
		"numbers": []int{1, 2, 3},
		"strings": []string{"a", "b"},
	}
	expected := map[string]interface{}{
		"numbers.[0]": 1,
		"numbers.[1]": 2,
		"numbers.[2]": 3,
		"strings.[0]": "a",
		"strings.[1]": "b",
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenComplexArray(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "John", "age": 30},
			map[string]interface{}{"name": "Jane", "age": 25},
		},
	}
	expected := map[string]interface{}{
		"users.[0].name": "John",
		"users.[0].age":  30,
		"users.[1].name": "Jane",
		"users.[1].age":  25,
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenWithPrefix(t *testing.T) {
	input := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	expected := map[string]interface{}{
		"user.name": "John",
		"user.age":  30,
	}
	result := Flatten(input, WithPrefix("user"))
	assert.Equal(t, expected, result)
}

func TestFlattenWithCustomSeparator(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}
	expected := map[string]interface{}{
		"user_name": "John",
		"user_age":  30,
	}
	result := Flatten(input, WithSep("_"))
	assert.Equal(t, expected, result)
}

func TestFlattenEmptyMap(t *testing.T) {
	input := map[string]interface{}{}
	expected := map[string]interface{}{}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenEmptyArray(t *testing.T) {
	input := map[string]interface{}{
		"empty": []interface{}{},
	}
	expected := map[string]interface{}{}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestExpandBasicFlat(t *testing.T) {
	input := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
	}
	expected := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestExpandNestedMap(t *testing.T) {
	input := map[string]interface{}{
		"user.name":         "John",
		"user.details.age":  30,
		"user.details.city": "NYC",
	}
	expected := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"details": map[string]interface{}{
				"age":  30,
				"city": "NYC",
			},
		},
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestExpandArray(t *testing.T) {
	input := map[string]interface{}{
		"numbers.[0]": 1,
		"numbers.[1]": 2,
		"numbers.[2]": 3,
		"strings.[0]": "a",
		"strings.[1]": "b",
	}
	expected := map[string]interface{}{
		"numbers": []interface{}{1, 2, 3},
		"strings": []interface{}{"a", "b"},
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestExpandComplexArray(t *testing.T) {
	input := map[string]interface{}{
		"users.[0].name": "John",
		"users.[0].age":  30,
		"users.[1].name": "Jane",
		"users.[1].age":  25,
	}
	expected := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "John", "age": 30},
			map[string]interface{}{"name": "Jane", "age": 25},
		},
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestExpandWithCustomSeparator(t *testing.T) {
	input := map[string]interface{}{
		"user_name": "John",
		"user_age":  30,
	}
	expected := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}
	result := Expand(input, WithSep("_"))
	assert.Equal(t, expected, result)
}

func TestExpandEmptyMap(t *testing.T) {
	input := map[string]interface{}{}
	result := Expand(input)
	assert.Nil(t, result)
}

func TestExpandSparseArray(t *testing.T) {
	input := map[string]interface{}{
		"items.[0]": "first",
		"items.[3]": "fourth",
	}
	expected := map[string]interface{}{
		"items": []interface{}{"first", nil, nil, "fourth"},
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestExpandMixedArrayAndMap(t *testing.T) {
	input := map[string]interface{}{
		"data.[0].name":     "John",
		"data.[0].tags.[0]": "admin",
		"data.[0].tags.[1]": "user",
		"metadata.version":  "1.0",
	}
	expected := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"name": "John",
				"tags": []interface{}{"admin", "user"},
			},
		},
		"metadata": map[string]interface{}{
			"version": "1.0",
		},
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestExpandNilValues(t *testing.T) {
	input := map[string]interface{}{
		"user.name":  "John",
		"user.email": nil,
		"user.age":   30,
	}
	expected := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "John",
			"email": nil,
			"age":   30,
		},
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestFlattenExpandRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"profile": map[string]interface{}{
						"name": "John",
						"age":  30,
					},
				},
			},
		},
		{
			name: "array",
			input: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
		},
		{
			name: "complex structure",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name":  "John",
						"roles": []string{"admin", "user"},
					},
				},
				"config": map[string]interface{}{
					"debug": true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flattened := Flatten(tt.input)
			expanded := Expand(flattened)

			originalJSON, _ := json.Marshal(tt.input)
			resultJSON, _ := json.Marshal(expanded)
			assert.Equal(t, string(originalJSON), string(resultJSON))
		})
	}
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
	src := A{
		F: 1,
		B: B{
			C: "test",
			D: 2,
		},
		Inner: Inner{},
	}
	actual := Flatten(src)
	assert.Equal(t, expecting, actual)
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

func TestFlattenMapWithNonStringKeys(t *testing.T) {
	// Maps with non-string keys should not be flattened
	input := map[int]string{
		1: "one",
		2: "two",
	}
	expected := map[string]interface{}{}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestFlattenInterface(t *testing.T) {
	var input interface{} = map[string]interface{}{
		"name": "John",
		"nested": map[string]interface{}{
			"age": 30,
		},
	}
	expected := map[string]interface{}{
		"name":       "John",
		"nested.age": 30,
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestExpandSingleValue(t *testing.T) {
	input := map[string]interface{}{
		"key": "value",
	}
	expected := map[string]interface{}{
		"key": "value",
	}
	result := Expand(input)
	assert.Equal(t, expected, result)
}

func TestExpandArrayWithGaps(t *testing.T) {
	input := map[string]interface{}{
		"items.[0]": "first",
		"items.[2]": "third",
		"items.[5]": "sixth",
	}
	result := Expand(input)
	resultMap := result.(map[string]interface{})
	items := resultMap["items"].([]interface{})

	assert.Equal(t, "first", items[0])
	assert.Equal(t, "third", items[2])
	assert.Equal(t, "sixth", items[5])
	assert.Equal(t, 6, len(items))
}

func TestFlattenWithBothPrefixAndSeparator(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}
	expected := map[string]interface{}{
		"root|user|name": "John",
		"root|user|age":  30,
	}
	result := Flatten(input, WithPrefix("root"), WithSep("|"))
	assert.Equal(t, expected, result)
}

func TestExpandWithMultipleSeparators(t *testing.T) {
	input := map[string]interface{}{
		"a|b|c": "value1",
		"x|y":   "value2",
	}
	expected := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value1",
			},
		},
		"x": map[string]interface{}{
			"y": "value2",
		},
	}
	result := Expand(input, WithSep("|"))
	assert.Equal(t, expected, result)
}

func TestFlattenDeepNesting(t *testing.T) {
	input := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": map[string]interface{}{
					"level4": map[string]interface{}{
						"value": "deep",
					},
				},
			},
		},
	}
	expected := map[string]interface{}{
		"level1.level2.level3.level4.value": "deep",
	}
	result := Flatten(input)
	assert.Equal(t, expected, result)
}

func TestExpandDeepNesting(t *testing.T) {
	input := map[string]interface{}{
		"level1.level2.level3.level4.value": "deep",
	}
	expected := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": map[string]interface{}{
					"level4": map[string]interface{}{
						"value": "deep",
					},
				},
			},
		},
	}
	result := Expand(input)
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
			privateNested: 123,
		},
	}

	// Test that flattening skips unexported fields and only includes exported ones
	expected := map[string]interface{}{
		"PublicField":               "public_value",
		"AnotherPublic":             42,
		"NestedStruct.PublicNested": "nested_public",
	}

	result := Flatten(input)
	assert.Equal(t, expected, result)

	// Verify that unexported fields are not included
	for key := range result {
		// All keys should start with uppercase (exported fields only)
		assert.True(t, key[0] >= 'A' && key[0] <= 'Z', "Found unexported field in result: %s", key)
	}

	// Verify that no unexported field names appear in the result
	for key := range result {
		assert.NotContains(t, key, "private", "Unexported field found in result")
		assert.NotContains(t, key, "another", "Unexported field found in result")
	}
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
		"PublicField":               "public_value",
		"AnotherPublic":             42,
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

func BenchmarkFlatten(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Flatten(example)
	}
}

func BenchmarkExpand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Expand(flat)
	}
}

func BenchmarkAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		flat := Flatten(example)
		_ = Expand(flat)
	}
}
