package fsmedia

import (
	"errors"
	"github.com/ottemo/foundation/db"
	"io/ioutil"
	"os"
)

// retrieve media storage name
func (it *FilesystemMediaStorage) GetName() string {
	return "FilesystemMediaStorage"
}

// returns path you can use to access media file (if possible for storage of course)
func (it *FilesystemMediaStorage) GetMediaPath(model string, objId string, mediaType string) (string, error) {
	return mediaType + "/" + model + "/" + objId + "/", nil
}

// retrieve contents of media entity for model object
func (it *FilesystemMediaStorage) Load(model string, objId string, mediaType string, mediaName string) ([]byte, error) {
	mediaPath, err := it.GetMediaPath(model, objId, mediaType)
	if err != nil {
		return nil, err
	}

	mediaFilePath := it.storageFolder + mediaPath + mediaName

	return ioutil.ReadFile(mediaFilePath)
}

// add media entity for model object
func (it *FilesystemMediaStorage) Save(model string, objId string, mediaType string, mediaName string, mediaData []byte) error {
	mediaPath, err := it.GetMediaPath(model, objId, mediaType)
	if err != nil {
		return err
	}

	mediaFolder := it.storageFolder + mediaPath
	mediaFilePath := mediaFolder + mediaName

	if _, err := os.Stat(mediaFolder); !os.IsExist(err) {
		err := os.MkdirAll(mediaFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	ioerr := ioutil.WriteFile(mediaFilePath, mediaData, os.ModePerm)
	if ioerr != nil {
		return ioerr
	}

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("Can't get database engine")
	}

	collection, err := dbEngine.GetCollection(MEDIA_DB_COLLECTION)
	if err != nil {
		return err
	}

	_, err = collection.Save(map[string]interface{}{"model": model, "object": objId, "type": mediaType, "media": mediaName})
	if err != nil {
		return err
	}

	return nil
}

// remove media entity for model object
func (it *FilesystemMediaStorage) Remove(model string, objId string, mediaType string, mediaName string) error {
	mediaPath, err := it.GetMediaPath(model, objId, mediaType)
	if err != nil {
		return err
	}

	mediaFilePath := it.storageFolder + mediaPath + mediaName

	// removing file
	err = os.Remove(mediaFilePath)
	if err != nil {
		return err
	}

	// removing DB records
	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("Can't get database engine")
	}

	collection, err := dbEngine.GetCollection(MEDIA_DB_COLLECTION)
	if err != nil {
		return err
	}

	err = collection.AddFilter("model", "=", model)
	err = collection.AddFilter("object", "=", objId)
	err = collection.AddFilter("type", "=", mediaType)
	err = collection.AddFilter("media", "=", mediaName)

	_, err = collection.Delete()
	return err
}

// get list of given type media entities for model object
func (it *FilesystemMediaStorage) ListMedia(model string, objId string, mediaType string) ([]string, error) {
	result := make([]string, 0)

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return result, errors.New("Can't get database engine")
	}

	collection, err := dbEngine.GetCollection(MEDIA_DB_COLLECTION)
	if err != nil {
		return result, err
	}

	collection.AddFilter("model", "=", model)
	collection.AddFilter("object", "=", objId)
	collection.AddFilter("type", "=", mediaType)

	records, err := collection.Load()

	for _, record := range records {
		result = append(result, record["media"].(string))
	}

	return result, err
}
