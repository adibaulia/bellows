package bellows

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// Test cases to achieve 100% coverage
func TestExpandArrayWithLargeGapsAndMixedTypes(t *testing.T) {
	// This test covers the case where we need to fill gaps with different types
	input := map[string]interface{}{
		"items.[0]":        "first",
		"items.[5]":        "sixth",
		"nested.[0].[2]":   "nested_value",
		"mixed.[0].name":   "john",
		"mixed.[3].age":    30,
	}
	result := Expand(input)
	resultMap := result.(map[string]interface{})
	
	// Check items array with gaps
	items := resultMap["items"].([]interface{})
	assert.Equal(t, "first", items[0])
	assert.Equal(t, "sixth", items[5])
	assert.Equal(t, 6, len(items))
	
	// Check nested array with gaps
	nested := resultMap["nested"].([]interface{})
	nestedInner := nested[0].([]interface{})
	assert.Equal(t, "nested_value", nestedInner[2])
	
	// Check mixed array with object gaps
	mixed := resultMap["mixed"].([]interface{})
	firstObj := mixed[0].(map[string]interface{})
	assert.Equal(t, "john", firstObj["name"])
	thirdObj := mixed[3].(map[string]interface{})
	assert.Equal(t, 30, thirdObj["age"])
}

func TestExpandInvalidArrayIndex(t *testing.T) {
	// This test covers the error handling in getArrayIndex
	input := map[string]interface{}{
		"items.[abc]": "invalid",
		"items.[]":    "empty",
		"normal.key":  "value",
	}
	result := Expand(input)
	resultMap := result.(map[string]interface{})
	
	// Invalid array indices should be treated as regular map keys
	items := resultMap["items"].(map[string]interface{})
	assert.Equal(t, "invalid", items["[abc]"])
	assert.Equal(t, "empty", items["[]"])
	
	normal := resultMap["normal"].(map[string]interface{})
	assert.Equal(t, "value", normal["key"])
}

func TestExpandArrayGapsWithArrayType(t *testing.T) {
	// This test specifically covers the case where we fill gaps with []interface{} type
	input := map[string]interface{}{
		"matrix.[0].[0]": "a",
		"matrix.[2].[1]": "b",
	}
	result := Expand(input)
	resultMap := result.(map[string]interface{})
	
	matrix := resultMap["matrix"].([]interface{})
	assert.Equal(t, 3, len(matrix))
	
	// First row
	firstRow := matrix[0].([]interface{})
	assert.Equal(t, "a", firstRow[0])
	
	// Second row should be empty array (gap filler)
	secondRow := matrix[1].([]interface{})
	assert.Equal(t, 0, len(secondRow))
	
	// Third row
	thirdRow := matrix[2].([]interface{})
	assert.Equal(t, "b", thirdRow[1])
}

func TestExpandArrayGapsWithMapType(t *testing.T) {
	// This test specifically covers the case where we fill gaps with map[string]interface{} type
	input := map[string]interface{}{
		"users.[0].name": "john",
		"users.[2].age":  30,
	}
	result := Expand(input)
	resultMap := result.(map[string]interface{})
	
	users := resultMap["users"].([]interface{})
	assert.Equal(t, 3, len(users))
	
	// First user
	firstUser := users[0].(map[string]interface{})
	assert.Equal(t, "john", firstUser["name"])
	
	// Second user should be empty map (gap filler)
	secondUser := users[1].(map[string]interface{})
	assert.Equal(t, 0, len(secondUser))
	
	// Third user
	thirdUser := users[2].(map[string]interface{})
	assert.Equal(t, 30, thirdUser["age"])
}

func TestExpandArrayIndexParseError(t *testing.T) {
	// This test covers the strconv.Atoi error case in getArrayIndex
	input := map[string]interface{}{
		"items.[1a2]": "invalid_number",
		"items.[3b]":  "another_invalid",
		"normal.key":  "value",
	}
	result := Expand(input)
	resultMap := result.(map[string]interface{})
	
	// Invalid numeric array indices should be treated as regular map keys
	items := resultMap["items"].(map[string]interface{})
	assert.Equal(t, "invalid_number", items["[1a2]"])
	assert.Equal(t, "another_invalid", items["[3b]"])
	
	normal := resultMap["normal"].(map[string]interface{})
	assert.Equal(t, "value", normal["key"])
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