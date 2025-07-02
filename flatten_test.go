package bellows

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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