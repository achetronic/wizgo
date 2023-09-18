package main

import (
	"encoding/json"
	"reflect"
)

func removeEmptyValuesFromJson(inputJSON []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(inputJSON, &data); err != nil {
		return []byte{}, err
	}

	result := removeEmpty(data)
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return []byte{}, err
	}

	return resultJSON, nil
}

func removeEmpty(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			if !isEmpty(val) {
				result[key] = removeEmpty(val)
			}
		}
		return result
	case []interface{}:
		var result []interface{}
		for _, val := range v {
			if !isEmpty(val) {
				result = append(result, removeEmpty(val))
			}
		}
		return result
	default:
		return v
	}
}

func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String, reflect.Map, reflect.Slice:
		return val.Len() == 0
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Float64:
		return false
	default:
		zero := reflect.Zero(val.Type())
		return reflect.DeepEqual(value, zero.Interface())
	}
}
