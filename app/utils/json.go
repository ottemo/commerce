package utils

import (
	"encoding/json"
	"errors"
)

// encodes inputData to JSON string if it's possible
func EncodeToJsonString(inputData interface{}) (string, error) {
	result, err := json.Marshal(inputData)
	return string(result), err
}

// decodes json string to map[string]interface{} if it's possible
func DecodeJsonToStringKeyMap(jsonData interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	var err error

	switch value := jsonData.(type) {
	case string:
		err = json.Unmarshal([]byte(value), &result)
	case []byte:
		err = json.Unmarshal(value, &result)
	default:
		err = errors.New("unsupported json data")
	}

	return result, err
}
