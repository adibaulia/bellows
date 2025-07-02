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

func TestAll(t *testing.T) {
	assert := assert.New(t)

	flattened := Flatten(example)
	expanded := Expand(flattened)

	expecting, _ := json.Marshal(example)
	actual, _ := json.Marshal(expanded)
	assert.Equal(expecting, actual)
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