package helper

import (
	"encoding/json"
)

func ModelToMap(model interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	modelByte, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(modelByte, &result); err != nil {
		return nil, err
	}

	return result, nil
}
