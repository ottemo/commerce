package product

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// AddMedia adds new media assigned to product
func (it *DefaultProduct) AddMedia(mediaType string, mediaName string, content []byte) error {
	productID := it.GetID()
	if productID == "" {
		return env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Save(it.GetModelName(), productID, mediaType, mediaName, content)
}

// RemoveMedia removes media assigned to product
func (it *DefaultProduct) RemoveMedia(mediaType string, mediaName string) error {
	productID := it.GetID()
	if productID == "" {
		return env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Remove(it.GetModelName(), productID, mediaType, mediaName)
}

// ListMedia lists media assigned to product
func (it *DefaultProduct) ListMedia(mediaType string) ([]string, error) {
	var result []string

	productID := it.GetID()
	if productID == "" {
		return result, env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMedia(it.GetModelName(), productID, mediaType)
}

// GetMedia returns content of media assigned to product
func (it *DefaultProduct) GetMedia(mediaType string, mediaName string) ([]byte, error) {
	productID := it.GetID()
	if productID == "" {
		return nil, env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.Load(it.GetModelName(), productID, mediaType, mediaName)
}

// GetMediaPath returns relative location of media assigned to product in media storage
func (it *DefaultProduct) GetMediaPath(mediaType string) (string, error) {
	productID := it.GetID()
	if productID == "" {
		return "", env.ErrorNew("product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), productID, mediaType)
}
