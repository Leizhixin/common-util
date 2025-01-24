package beanutil

import "strconv"

func ParseResultToMap(raw map[interface{}]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for k, v := range raw {
		//resultMap[k.(string)] = v
		var newKey = ""
		var newValue interface{}
		switch k.(type) {
		case string:
			newKey = k.(string)
		case int:
			newKey = strconv.Itoa(k.(int))
		}
		if len(newKey) == 0 || newKey == "class" {
			continue
		}
		switch v.(type) {
		case map[interface{}]interface{}:
			newValue = ParseResultToMap(v.(map[interface{}]interface{}))
		case []interface{}:
			HandleList(v.([]interface{}))
			newValue = v
		default:
			newValue = v
		}
		resultMap[newKey] = newValue
	}
	return resultMap
}

func HandleList(rawList []interface{}) {
	for i, raw := range rawList {
		switch raw.(type) {
		case map[interface{}]interface{}:
			rawList[i] = ParseResultToMap(raw.(map[interface{}]interface{}))
		case []interface{}:
			HandleList(rawList[i].([]interface{}))
		}
	}
}
