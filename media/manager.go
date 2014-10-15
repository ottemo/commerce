package media

import (
	"github.com/ottemo/foundation/env"
)

var currentMediaStorage IMediaStorage = nil
var callbacksOnMediaStorageStart = []func() error{}

func RegisterOnMediaStorageStart(callback func() error) {
	callbacksOnMediaStorageStart = append(callbacksOnMediaStorageStart, callback)
}

func OnMediaStorageStart() error {
	for _, callback := range callbacksOnMediaStorageStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

func RegisterMediaStorage(newEngine IMediaStorage) error {
	if currentMediaStorage == nil {
		currentMediaStorage = newEngine
	} else {
		return env.ErrorNew("Sorry, '" + currentMediaStorage.GetName() + "' media storage already registered")
	}
	return nil
}

func GetMediaStorage() (IMediaStorage, error) {
	if currentMediaStorage != nil {
		return currentMediaStorage, nil
	} else {
		return nil, env.ErrorNew("no registered media storage")
	}
}
