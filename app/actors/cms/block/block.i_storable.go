package block

import (
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
)

// returns id for cms block
func (it *DefaultCMSBlock) GetId() string {
	return it.id
}

// sets id for cms block
func (it *DefaultCMSBlock) SetId(newId string) error {
	it.id = newId
	return nil
}

// loads cms block information from DB
func (it *DefaultCMSBlock) Load(id string) error {
	collection, err := db.GetCollection(CMS_BLOCK_COLLECTION_NAME)
	if err != nil {
		return err
	}

	dbValues, err := collection.LoadById(id)
	if err != nil {
		return err
	}

	it.SetId(utils.InterfaceToString(dbValues["_id"]))

	it.Content = utils.InterfaceToString(dbValues["content"])
	it.Identifier = utils.InterfaceToString(dbValues["identifier"])
	it.CreatedAt = utils.InterfaceToTime(dbValues["created_at"])
	it.UpdatedAt = utils.InterfaceToTime(dbValues["updated_at"])

	return nil
}

// removes current cms block from DB
func (it *DefaultCMSBlock) Delete() error {
	collection, err := db.GetCollection(CMS_BLOCK_COLLECTION_NAME)
	if err != nil {
		return err
	}

	err = collection.DeleteById(it.GetId())
	if err != nil {
		return err
	}

	return err
}

// stores current cms block to DB
func (it *DefaultCMSBlock) Save() error {
	collection, err := db.GetCollection(CMS_BLOCK_COLLECTION_NAME)
	if err != nil {
		return err
	}

	// packing data before save
	storingValues := make(map[string]interface{})

	storingValues["_id"] = it.GetId()

	storingValues["identifier"] = it.GetIdentifier()
	storingValues["content"] = it.GetContent()

	currentTime := time.Now()

	if it.CreatedAt.IsZero() {
		storingValues["created_at"] = currentTime
	}
	storingValues["updated_at"] = currentTime

	newId, err := collection.Save(storingValues)
	if err != nil {
		return err
	}
	it.SetId(newId)

	return nil
}
