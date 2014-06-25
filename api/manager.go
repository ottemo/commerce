package api

import (
	"errors"
)

var currentService IService = nil
var callbacksOnServiceStart = []func() error{}

func RegisterOnServiceStart(callback func() error) {
	callbacksOnServiceStart = append(callbacksOnServiceStart, callback)
}
func OnServiceStart() error {
	for _, callback := range callbacksOnServiceStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

func RegisterService(newService IService) error {
	if currentService == nil {
		currentService = newService
	} else {
		return errors.New("Sorry, '" + currentService.GetName() + "' already registered")
	}
	return nil
}

func GetService() IService {
	return currentService
}
