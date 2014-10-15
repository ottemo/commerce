package db

import (
	"github.com/ottemo/foundation/env"
)

var currentDBEngine I_DBEngine = nil

var callbacksOnDatabaseStart = []func() error{}

func RegisterOnDatabaseStart(callback func() error) {
	callbacksOnDatabaseStart = append(callbacksOnDatabaseStart, callback)
}
func OnDatabaseStart() error {
	for _, callback := range callbacksOnDatabaseStart {
		if err := callback(); err != nil {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

func RegisterDBEngine(newEngine I_DBEngine) error {
	if currentDBEngine == nil {
		currentDBEngine = newEngine
	} else {
		return env.ErrorNew("Sorry, '" + currentDBEngine.GetName() + "' already registered")
	}
	return nil
}

func GetDBEngine() I_DBEngine {
	return currentDBEngine
}
