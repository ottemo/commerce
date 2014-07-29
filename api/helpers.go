package api

import (
	"errors"
)



// tries to represent HTTP request content in map[string]interface{} format
func GetRequestContentAsMap(params *T_APIHandlerParams) (map[string]interface{}, error) {

	result, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, errors.New("unexpected request content")
		} else {
			result = make(map[string]interface{})
		}
	}

	return result, nil
}
