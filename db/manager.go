package db

import (
	"github.com/ottemo/foundation/env"
)

var (
	currentDBEngine          I_DBEngine = nil              // currently registered database service in system
	callbacksOnDatabaseStart            = []func() error{} // set of callback function on database service start
)

// registers new callback on database service start
func RegisterOnDatabaseStart(callback func() error) {
	callbacksOnDatabaseStart = append(callbacksOnDatabaseStart, callback)
}

// fires database service start event (callback handling)
func OnDatabaseStart() error {
	for _, callback := range callbacksOnDatabaseStart {
		if err := callback(); err != nil {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// registers database service in the system
//   - will cause error if there are couple candidates for that role
func RegisterDBEngine(newEngine I_DBEngine) error {
	if currentDBEngine == nil {
		currentDBEngine = newEngine
	} else {
		return env.ErrorNew("Sorry, '" + currentDBEngine.GetName() + "' already registered")
	}
	return nil
}

// returns currently used database service implementation
func GetDBEngine() I_DBEngine {
	return currentDBEngine
}
