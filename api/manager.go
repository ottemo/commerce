package api

import (
	"github.com/ottemo/foundation/env"
)

var currentRestService I_RestService = nil
var callbacksOnRestServiceStart = []func() error{}

func RegisterOnRestServiceStart(callback func() error) {
	callbacksOnRestServiceStart = append(callbacksOnRestServiceStart, callback)
}
func OnRestServiceStart() error {
	for _, callback := range callbacksOnRestServiceStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

func RegisterRestService(newService I_RestService) error {
	if currentRestService == nil {
		currentRestService = newService
	} else {
		return env.ErrorNew("Sorry, '" + currentRestService.GetName() + "' already registered")
	}
	return nil
}

func GetRestService() I_RestService {
	return currentRestService
}
