package fsmedia

import (
	"errors"
	"github.com/ottemo/foundation/database"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/config"
)

func init() {
	instance := new( FilesystemMediaStorage )
	if err := media.RegisterMediaStorage( instance ); err == nil {
		instance.setupWaitCnt = 2

		config.RegisterOnConfigIniStart(instance.setupOnIniConfig)
		database.RegisterOnDatabaseStart(instance.setupOnDatabase)
	}
}

func (it *FilesystemMediaStorage) setupCheckDone() {
	if it.setupWaitCnt--; it.setupWaitCnt == 0 {
		media.OnMediaStorageStart()
	}
}

func (it *FilesystemMediaStorage) setupOnIniConfig() error {

	var storageFolder = MEDIA_DEFAULT_FOLDER

	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("media.fsmedia.folder"); iniValue != "" {
			storageFolder = iniValue
		}
	}

	// TODO: add checks for folder existence and rights
	it.storageFolder = storageFolder

	it.setupCheckDone()

	return nil
}

func (it *FilesystemMediaStorage) setupOnDatabase() error {

	dbEngine := database.GetDBEngine()
	if dbEngine == nil { return errors.New("Can't get database engine") }

	collection, err := dbEngine.GetCollection( MEDIA_DB_COLLECTION )
	if err != nil { return err }

	// TODO: make 3 column PK constraint (model, object, media)
	collection.AddColumn("model", "text", true)
	collection.AddColumn("object", "text", true)
	collection.AddColumn("type", "text", true)
	collection.AddColumn("media", "text", false)

	it.setupCheckDone()

	return nil
}
