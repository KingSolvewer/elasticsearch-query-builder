package elastic

import (
	"errors"
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"reflect"
)

func CheckHitsDestType(dest any) error {
	switch dest.(type) {
	case nil, map[string]any, *map[string]any, *[]map[string]any:
		return nil
	default:
		reflectPtr := reflect.ValueOf(dest)
		if reflectPtr.Kind() != reflect.Ptr {
			return errors.New("expected Pointer dest")
		}

		reflectValue := reflectPtr.Elem()
		switch reflectValue.Kind() {
		case reflect.Array, reflect.Slice, reflect.Struct:
			return nil
		default:
			return errors.New("only supports parsing variables of the types map, *map, *struct, *[]struct, *[]map")
		}
	}
}

func CheckTopHitsDestType(dest any) error {
	switch dest.(type) {
	case nil, map[string]any, *map[string]any:
		return nil
	default:
		reflectPtr := reflect.ValueOf(dest)
		if reflectPtr.Kind() != reflect.Ptr {
			return errors.New("expected Pointer dest")
		}

		reflectStruct := reflectPtr.Elem()
		if reflectStruct.Kind() != reflect.Struct {
			return errors.New("the element that the Pointer dest points to is not struct")
		}

		return nil
	}
}

func checkType(value any) bool {
	switch value.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, string, bool:
		return true
	default:
		return false
	}
}

func checkAggsRangeType(ranges []aggs.Ranges) bool {
	var check bool
	for _, r := range ranges {
		switch r.From.(type) {
		case int, string:
			check = true
		}
		switch r.To.(type) {
		case int, string:
			check = true
		}
	}

	return check
}
