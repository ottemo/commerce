package fsmedia

import (
	"bytes"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"
)

// GetName returns media storage name
func (it *FilesystemMediaStorage) GetName() string {
	return "FilesystemMediaStorage"
}

// GetMediaPath returns path you can use to access media file (if possible for storage of course)
func (it *FilesystemMediaStorage) GetMediaPath(model string, objID string, mediaType string) (string, error) {
	return mediaType + "/" + model + "/" + objID + "/", nil
}

// Load retrieves contents of media entity for given model object
func (it *FilesystemMediaStorage) Load(model string, objID string, mediaType string, mediaName string) ([]byte, error) {
	mediaPath, err := it.GetMediaPath(model, objID, mediaType)
	if err != nil {
		return nil, err
	}

	mediaFilePath := it.storageFolder + mediaPath + mediaName

	return ioutil.ReadFile(mediaFilePath)
}

// Save adds media entity for model object
func (it *FilesystemMediaStorage) Save(model string, objID string, mediaType string, mediaName string, mediaData []byte) error {
	mediaPath, err := it.GetMediaPath(model, objID, mediaType)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	mediaFolder := it.storageFolder + mediaPath
	mediaFilePath := mediaFolder + mediaName

	if _, err := os.Stat(mediaFolder); !os.IsExist(err) {
		err := os.MkdirAll(mediaFolder, os.ModePerm)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	ioerr := ioutil.WriteFile(mediaFilePath, mediaData, os.ModePerm)
	if ioerr != nil {
		return env.ErrorDispatch(ioerr)
	}

	// we have image associated media, so making special treatment
	if mediaType == ConstMediaTypeImage {

		// checking that image is png or jpeg, making it jpeg if not
		decodedImage, imageFormat, err := image.Decode(bytes.NewReader(mediaData))
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// converting image to known format if needed
		if imageFormat != "jpeg" && imageFormat != "png" {

			newFile, err := os.Create(mediaFilePath)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			idx := strings.LastIndex(mediaName, ".")
			if idx != -1 {
				mediaFilePath = mediaName[idx:]
			}
			mediaFilePath += ".jpg"

			err = jpeg.Encode(newFile, decodedImage, nil)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}

		// resizing to image sizes system currently using
		for imageSize := range it.imageSizes {
			err = it.ResizeMediaImage(model, objID, mediaName, imageSize)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

	// making database record
	//------------------------

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("Can't get database engine")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection.AddFilter("model", "=", model)
	dbCollection.AddFilter("object", "=", objID)
	dbCollection.AddFilter("type", "=", mediaType)
	dbCollection.AddFilter("media", "=", mediaName)

	count, err := dbCollection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if count == 0 {
		_, err = dbCollection.Save(map[string]interface{}{"model": model, "object": objID, "type": mediaType, "media": mediaName})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// Remove removes media entity for model object
func (it *FilesystemMediaStorage) Remove(model string, objID string, mediaType string, mediaName string) error {

	// preparing DB collection
	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("Can't get database engine")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("model", "=", model)
	err = dbCollection.AddFilter("object", "=", objID)
	err = dbCollection.AddFilter("type", "=", mediaType)
	err = dbCollection.AddFilter("media", "=", mediaName)

	// removing files
	records, err := dbCollection.Load()

	mediaFolder := it.storageFolder

	for _, record := range records {
		if mediaName, ok := record["media"].(string); ok {

			if path, err := it.GetMediaPath(model, objID, mediaType); err == nil {

				// looking for object image sizes to remove
				if mediaType == ConstMediaTypeImage {
					for imageSize := range it.imageSizes {
						mediaFilePath := mediaFolder + path + it.GetResizedMediaName(mediaName, imageSize)
						os.Remove(mediaFilePath)
					}
				}

				os.Remove(mediaFolder + path + mediaName)
			}
		}
	}

	// removing DB records
	_, err = dbCollection.Delete()

	return env.ErrorDispatch(err)
}

// ListMedia returns list of given type media entities for a given model object
func (it *FilesystemMediaStorage) ListMedia(model string, objID string, mediaType string) ([]string, error) {
	var result []string

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return result, env.ErrorNew("Can't get database engine")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	dbCollection.AddFilter("model", "=", model)
	dbCollection.AddFilter("object", "=", objID)
	dbCollection.AddFilter("type", "=", mediaType)

	records, err := dbCollection.Load()

	for _, record := range records {
		if mediaName, ok := record["media"].(string); ok {
			result = append(result, mediaName)
		}
	}

	// checking that object have all image sizes
	if mediaType == ConstMediaTypeImage {
		if path, err := it.GetMediaPath(model, objID, mediaType); err == nil {
			for _, mediaName := range result {
				for imageSize := range it.imageSizes {
					mediaFilePath := path + it.GetResizedMediaName(mediaName, imageSize)

					if _, err := os.Stat(mediaFilePath); os.IsNotExist(err) {
						it.ResizeMediaImage(model, objID, mediaName, imageSize)
					}
				}
			}
		}
	}

	return result, env.ErrorDispatch(err)
}
