package db

import (
	"errors"
)

var currentDBEngine I_DBEngine = nil

var callbacksOnDatabaseStart = []func() error {}

func RegisterOnDatabaseStart(callback func() error) {
	callbacksOnDatabaseStart = append(callbacksOnDatabaseStart, callback)
}
func OnDatabaseStart() error {
	for _, callback := range callbacksOnDatabaseStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

func RegisterDBEngine(newEngine I_DBEngine) error {
	if currentDBEngine == nil {
		currentDBEngine = newEngine
	} else {
		return errors.New("Sorry, '" + currentDBEngine.GetName() + "' already registered")
	}
	return nil
}

func GetDBEngine() I_DBEngine {
	return currentDBEngine
}
