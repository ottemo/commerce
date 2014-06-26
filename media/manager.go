package media

import (
	"errors"
)

var currentMediaStorage IMediaStorage = nil
var callbacksOnMediaStorageStart = []func() error {}

func RegisterOnMediaStorageStart(callback func() error) {
	callbacksOnMediaStorageStart = append(callbacksOnMediaStorageStart, callback)
}

func OnMediaStorageStart () error {
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
		return errors.New("Sorry, '" + currentMediaStorage.GetName() + "' media storage already registered")
	}
	return nil
}

func GetMediaStorage() IMediaStorage {
	return currentMediaStorage
}
