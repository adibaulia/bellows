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

	flat := Flatten(example)
	expanded := Expand(flat)

	expecting, _ := json.Marshal(example)
	actual, _ := json.Marshal(expanded)
	assert.Equal(expecting, actual)
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
