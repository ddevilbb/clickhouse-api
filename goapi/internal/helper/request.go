package helper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ParseQueryValues(queryValues url.Values) (map[string][]string, error) {
	result := make(map[string][]string)
	if len(queryValues) > 0 {
		var filteredKey, value string
		regEx := regexp.MustCompile(`(?m)(\w+)\[\d*\]`)
		for key, values := range queryValues {
			filteredKey = regEx.ReplaceAllString(key, "$1")
			value = strings.Join(values, ", ")
			if len(value) == 0 {
				continue
			}
			if result[filteredKey] == nil {
				result[filteredKey] = []string{}
			}
			result[filteredKey] = append(result[filteredKey], values...)
		}
	}

	return result, nil
}

func GetModelQueryValues(queryParams map[string][]string, modelMap map[string]interface{}) (map[string][]string, error) {
	result := make(map[string][]string)
	if len(queryParams) > 0 {
		var filteredKey string
		regEx := regexp.MustCompile(`(?m)(\w+)\[\w+\d*\]`)
		for key, values := range queryParams {
			filteredKey = regEx.ReplaceAllString(key, "$1")
			if filteredKey == "order" || filteredKey == "page" || filteredKey == "limit" {
				continue
			}
			if _, ok := modelMap[key]; !ok {
				return nil, fmt.Errorf("field %s not exists in Model %+v", key, modelMap)
			}
			if result[key] == nil {
				result[key] = []string{}
			}
			result[key] = append(result[key], values...)
		}
	}

	return result, nil
}

func GetOrderQueryValues(queryParams map[string][]string, modelMap map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string)
	if len(queryParams) > 0 {
		var filteredKey, orderKey string
		regEx := regexp.MustCompile(`(?m)(\w+)\[(\w*)\]`)
		for key, value := range queryParams {
			filteredKey = regEx.ReplaceAllString(key, "$1")
			orderKey = regEx.ReplaceAllString(key, "$2")
			if filteredKey != "order" {
				continue
			}
			if _, ok := modelMap[orderKey]; !ok {
				return nil, fmt.Errorf("field %s not exists in Model %+v", key, modelMap)
			}
			result[orderKey] = value[0]
		}
	}

	return result, nil
}
