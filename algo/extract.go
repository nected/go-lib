package algo

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/nected/go-lib/utils"
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
		} else if reflect.TypeOf(source).Kind() == reflect.Map {
			ok := false
			//converting source into map
			var sourceMap map[string]interface{}
			err := utils.JsonToStruct(source, &sourceMap)
			if err != nil {
				return nil, fmt.Errorf("item is not map %v", source)
			}

			key, indexes, err := extractKeyIndex(itemKeys[i])
			if err != nil {
				return nil, err
			}
			// retrieving item using key
			source, ok = sourceMap[key]
			if !ok {
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
		} else if reflect.TypeOf(source).Kind() == reflect.Slice ||
			reflect.TypeOf(source).Kind() == reflect.Array ||
			reflect.TypeOf(source).Kind() == reflect.String {
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
		} else {
			return nil, fmt.Errorf("inavlid usage of %v key non map/list", keyStr)
		}
	}
	return source, nil
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

	// check if last key is list with index
	for _, index := range indexes {
		if reflect.TypeOf(source).Kind() == reflect.Slice ||
			reflect.TypeOf(source).Kind() == reflect.Array {
			var sourceArr []interface{}
			err := utils.JsonToStruct(source, &sourceArr)
			if err != nil {
				return nil, fmt.Errorf("item is not list %v", err)
			}

			// if list index doesn't exist
			if len(sourceArr) <= index {
				if missingKeyError {
					return nil, fmt.Errorf("index out of bound %v", index)
				}
				return nil, nil
			}
			// retrieving item using index
			source = sourceArr[index]
		} else if reflect.TypeOf(source).Kind() == reflect.String {
			sourceStr := source.(string)
			if len(sourceStr) <= index {
				if missingKeyError {
					return nil, fmt.Errorf("index out of bound %v", index)
				}
				return nil, nil
			}
			// retrieving item using index
			source = string(sourceStr[index])
		} else {
			return nil, fmt.Errorf("inavalid usage of index")
		}
	}
	return source, nil
}
