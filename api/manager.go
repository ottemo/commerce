package api

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	currentRestService          InterfaceRestService = nil              // currently registered RESTFul service in system
	callbacksOnRestServiceStart                      = []func() error{} // set of callback function on RESTFul service start
)

// registers new callback on RESTFul service start
func RegisterOnRestServiceStart(callback func() error) {
	callbacksOnRestServiceStart = append(callbacksOnRestServiceStart, callback)
}

// fires RESTFul service start event (callback handling)
func OnRestServiceStart() error {
	for _, callback := range callbacksOnRestServiceStart {
		if err := callback(); err != nil {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// registers RESTFul service in the system
//   - will cause error if there are couple candidates for that role
func RegisterRestService(newService InterfaceRestService) error {
	if currentRestService == nil {
		currentRestService = newService
	} else {
		return env.ErrorNew("Sorry, '" + currentRestService.GetName() + "' already registered")
	}
	return nil
}

// returns currently used RESTFul service implementation
func GetRestService() InterfaceRestService {
	return currentRestService
}
