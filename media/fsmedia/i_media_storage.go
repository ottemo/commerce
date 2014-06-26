package fsmedia

import (
	"os"
	"errors"
	"io/ioutil"
	"github.com/ottemo/foundation/database"
)

func (it *FilesystemMediaStorage) GetName() string {
	return "FilesystemMediaStorage"
}

func (it *FilesystemMediaStorage) Load(model string, objId string, mediaType string, mediaName string) ([]byte, error) {
	fullFileName := it.storageFolder + "/" + mediaType + "/" + model + "/" + mediaName

	return ioutil.ReadFile(fullFileName)
}

func (it *FilesystemMediaStorage) Save(model string, objId string, mediaType string, mediaName string, mediaData []byte) error {

	filePath := it.storageFolder + "/" + mediaType + "/" + model
	fullFileName := filePath + "/" + mediaName

	if _, err := os.Stat(filePath); !os.IsExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil { return err }
	}

	ioerr := ioutil.WriteFile(fullFileName, mediaData, os.ModePerm)
	if ioerr != nil { return ioerr }


	dbEngine := database.GetDBEngine()
	if dbEngine == nil { return errors.New("Can't get database engine") }

	collection, err := dbEngine.GetCollection( MEDIA_DB_COLLECTION )
	if err != nil { return err }

	_, err = collection.Save( map[string]interface{} { "model": model, "object": objId, "type": mediaType, "media": mediaName } )
	if err != nil { return err }

	return nil
}

func (it *FilesystemMediaStorage) ListMedia(model string, objId string, mediaType string) ([]string, error) {
	result := make([]string, 0)

	dbEngine := database.GetDBEngine()
	if dbEngine == nil { return nil, errors.New("Can't get database engine") }

	collection, err := dbEngine.GetCollection( MEDIA_DB_COLLECTION )
	if err != nil { return nil, err }

	collection.AddFilter("model", "=", model)
	collection.AddFilter("object", "=", objId)
	collection.AddFilter("type", "=", mediaType)

	collection.Load()

	return result, nil
}





