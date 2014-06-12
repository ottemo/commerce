package app

import (
	"github.com/ottemo/foundation/rest_service"
)

var callbacksOnAppStart = []func() error{}

// OnAppStart is a method to register callbacks upon initialization of Foundation.
func OnAppStart(callback func() error) {
	callbacksOnAppStart = append(callbacksOnAppStart, callback)
}

// Start iterates through the registered callbacks with Foundation is first started.
func Start() error {
	for _, callback := range callbacksOnAppStart {
		if err := callback(); err != nil {
			return err
		}
	}

	return nil
}

// Serve starts and returns the REST Endpoint.
func Serve() error {
	return rest_service.GetRestService().Run()
}
