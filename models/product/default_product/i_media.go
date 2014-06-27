package default_product

import (
	"errors"
	"github.com/ottemo/foundation/media"
)

func (it *DefaultProductModel) AddMedia( mediaType string, mediaName string, content []byte ) error {
	productId := it.GetId()
	if productId == "" { return errors.New("product id not set") }

	mediaStorage, err := media.GetMediaStorage()
	if err != nil { return err }

	return mediaStorage.Save(it.GetModelName(), productId, mediaType, mediaName, content)
}

func (it *DefaultProductModel) RemoveMedia( mediaType string, mediaName string ) error {
	productId := it.GetId()
	if productId == "" { return errors.New("product id not set") }

	mediaStorage, err := media.GetMediaStorage()
	if err != nil { return err }

	return 	mediaStorage.Remove(it.GetModelName(), productId, mediaType, mediaName)
}

func (it *DefaultProductModel) ListMedia( mediaType string ) ([]string, error) {
	result := make([]string, 0)

	productId := it.GetId()
	if productId == "" { return result, errors.New("product id not set") }

	mediaStorage, err := media.GetMediaStorage()
	if err != nil { return result, err }

	return mediaStorage.ListMedia(it.GetModelName(), productId, mediaType)
}

func (it *DefaultProductModel) GetMedia( mediaType string, mediaName string ) ([]byte, error) {
	productId := it.GetId()
	if productId == "" { return nil, errors.New("product id not set") }

	mediaStorage, err := media.GetMediaStorage()
	if err != nil { return nil, err }

	return mediaStorage.Load(it.GetModelName(), productId, mediaType, mediaName)
}
