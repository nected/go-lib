package algo

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	ERROR_MISSING_KEY_VALUE = "error"
)

var listIndexRegexMatcher = regexp.MustCompile(`^([\w-]*)(\[([0-9]+)\])+$`)

// common function to extract any key from source
// example: source -> map, key-> k[0][1].k1
// example source -> list, key -> [0][1]
// example source -> string, key [0]
func GetValFromSource(source interface{}, keyStr string, options ...string) (interface{}, error) {
	if keyStr == "" {
		return source, nil
	}

	// default missing key behaviour is null value
	var MISING_KEY_ERROR bool
	if len(options) > 0 && options[0] == ERROR_MISSING_KEY_VALUE {
		MISING_KEY_ERROR = true
	}

	itemKeys := strings.Split(keyStr, ".")
	for i := 0; i < len(itemKeys); i++ {
		if source == nil {
			if MISING_KEY_ERROR {
				return nil, fmt.Errorf("source is null for key %v", itemKeys[i])
			}
			return nil, nil
		}
		switch reflect.TypeOf(source).Kind() {
		case reflect.Map, reflect.Struct:
			key, indexes, err := extractKeyIndex(itemKeys[i])
			if err != nil {
				return nil, err
			}

			// retrieving item using key
			source, err = getMapOrStructKeyValue(source, key)
			if err != nil {
				if MISING_KEY_ERROR {
					return nil, fmt.Errorf("key %v is not present", itemKeys[i])
				}
				return nil, nil
			}
			// retrieving index if present in key
			source, err = getListIndexValue(source, indexes, MISING_KEY_ERROR)
			if err != nil {
				return nil, err
			}
		case reflect.Slice, reflect.Array, reflect.String:
			key, indexes, err := extractKeyIndex(itemKeys[i])
			if err != nil {
				return nil, err
			}
			if key != "" {
				return nil, fmt.Errorf("key used on list item %v", keyStr)
			}
			// retrieving index if present in key
			source, err = getListIndexValue(source, indexes, MISING_KEY_ERROR)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("inavlid usage of %v key non map/list", keyStr)
		}
	}
	return source, nil
}
func getMapOrStructKeyValue(m interface{}, key string) (interface{}, error) {
	if m == nil {
		return nil, nil
	}

	switch v := reflect.ValueOf(m); v.Kind() {
	case reflect.Map:
		keyValue := reflect.ValueOf(key)
		mapValue := v.MapIndex(keyValue)
		if !mapValue.IsValid() {
			return nil, fmt.Errorf("inavlid value: %v", mapValue)
		}
		return mapValue.Interface(), nil
	case reflect.Struct:
		f := v.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, key)
		})
		if !f.IsValid() {
			return nil, fmt.Errorf("inavlid value for: %v", key)
		}

		if f.Kind() == reflect.Ptr && !f.IsNil() {
			f = f.Elem()
		}
		return f.Interface(), nil
	default:
		return nil, fmt.Errorf("%v is not a map", v.Kind())
	}
}

func getArrayIndexValue(arr any, idx int, missingKeyError bool) (any, error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("%v is not a array", v.Kind())
	}
	if v.Len() <= idx {
		if missingKeyError {
			return nil, fmt.Errorf("out_of_index")
		}
		return nil, nil
	}

	idxVal := v.Index(idx)
	if !idxVal.IsValid() {
		return nil, fmt.Errorf("inavlid value: %v", idxVal)
	}
	return idxVal.Interface(), nil
}

func extractKeyIndex(key string) (string, []int, error) {
	indexes := make([]int, 0)
	key = strings.ReplaceAll(key, " ", "")
	//matching if index used on with or without key
	// matched k[0][0] or [0][0]
	reg := listIndexRegexMatcher.FindAllString(key, -1)
	if len(reg) > 0 {
		keyArr := strings.Split(strings.ReplaceAll(key, "]", ""), "[")
		if len(keyArr) > 1 {
			//  building index
			for j := 1; j < len(keyArr); j++ {
				ind, err := strconv.Atoi(keyArr[j])
				if err != nil {
					return "", nil, fmt.Errorf("inavlid index in the key %v", key)
				}
				indexes = append(indexes, ind)
			}
			// overriding key if indexes persent
			key = keyArr[0]
		}
	}
	return key, indexes, nil
}

func getListIndexValue(source interface{}, indexes []int, missingKeyError bool) (interface{}, error) {
	if len(indexes) > 0 && source == nil {
		return nil, fmt.Errorf("source must not be nil")
	}

	var err error
	// check if last key is list with index
	for _, index := range indexes {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Slice, reflect.Array:
			source, err = getArrayIndexValue(source, index, missingKeyError)
			if err != nil {
				return nil, err
			}
		case reflect.String:
			sourceStr := source.(string)
			if len(sourceStr) <= index {
				if missingKeyError {
					return nil, fmt.Errorf("index out of bound %v", index)
				}
				return nil, nil
			}
			// retrieving item using index
			source = string(sourceStr[index])
		default:
			return nil, fmt.Errorf("inavalid usage of index")
		}
	}
	return source, nil
}

func getField(obj interface{}, name string) (interface{}, bool) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	f := v.FieldByNameFunc(func(s string) bool {
		return strings.EqualFold(s, name)
	})
	if !f.IsValid() {
		return nil, false
	}
	return f.Interface(), true
}
