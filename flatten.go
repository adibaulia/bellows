package bellows

import (
	"fmt"
	"reflect"
)

func Flatten(value interface{}, sep string) map[string]interface{} {
	return FlattenPrefixed(value, sep, "")
}

func FlattenPrefixed(value interface{}, sep string, prefix string) map[string]interface{} {
	m := make(map[string]interface{}, 5)
	FlattenPrefixedToResult(value, sep, prefix, m)
	return m
}

func FlattenPrefixedToResult(value interface{}, sep string, prefix string, m map[string]interface{}) {
	original := reflect.ValueOf(value)
	kind := original.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		original = reflect.Indirect(original)
		kind = original.Kind()
	}

	if !original.IsValid() {
		if prefix != "" {
			m[prefix] = nil
		}
		return
	}

	t := original.Type()

	switch kind {
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			break
		}
		keys := original.MapKeys()
		base := ""
		if prefix != "" {
			base = prefix + sep
		}
		for _, childKey := range keys {
			childValue := original.MapIndex(childKey)
			FlattenPrefixedToResult(childValue.Interface(), sep, base+childKey.String(), m)
		}
	case reflect.Struct:
		numField := original.NumField()
		base := ""
		for i := 0; i < numField; i++ {
			childValue := original.Field(i)
			f := t.Field(i)
			childKey := f.Name
			if f.Anonymous {
				childKey = ""
				base = prefix
			} else if prefix != "" {
				base = prefix + sep
			}
			FlattenPrefixedToResult(childValue.Interface(), sep, base+childKey, m)
		}
	case reflect.Array, reflect.Slice:
		l := original.Len()
		base := prefix
		for i := 0; i < l; i++ {
			childValue := original.Index(i)
			FlattenPrefixedToResult(childValue.Interface(), sep, fmt.Sprintf("%s%s[%d]", base, sep, i), m)
		}
	default:
		if prefix != "" {
			m[prefix] = value
		}
	}
}
