package product

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// adds new media assigned to product
func (it *DefaultProduct) AddMedia(mediaType string, mediaName string, content []byte) error {
	productId := it.GetId()
	if productId == "" {
		return env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Save(it.GetModelName(), productId, mediaType, mediaName, content)
}

// removes media assigned to product
func (it *DefaultProduct) RemoveMedia(mediaType string, mediaName string) error {
	productId := it.GetId()
	if productId == "" {
		return env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Remove(it.GetModelName(), productId, mediaType, mediaName)
}

// lists media assigned to product
func (it *DefaultProduct) ListMedia(mediaType string) ([]string, error) {
	result := make([]string, 0)

	productId := it.GetId()
	if productId == "" {
		return result, env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMedia(it.GetModelName(), productId, mediaType)
}

// returns content of media assigned to product
func (it *DefaultProduct) GetMedia(mediaType string, mediaName string) ([]byte, error) {
	productId := it.GetId()
	if productId == "" {
		return nil, env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.Load(it.GetModelName(), productId, mediaType, mediaName)
}

// returns relative location of media assigned to product in media storage
func (it *DefaultProduct) GetMediaPath(mediaType string) (string, error) {
	productId := it.GetId()
	if productId == "" {
		return "", env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), productId, mediaType)
}
