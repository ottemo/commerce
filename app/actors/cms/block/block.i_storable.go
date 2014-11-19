package block

import (
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
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
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbValues, err := collection.LoadById(id)
	if err != nil {
		return env.ErrorDispatch(err)
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
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteById(it.GetId())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return env.ErrorDispatch(err)
}

// stores current cms block to DB
func (it *DefaultCMSBlock) Save() error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
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
		return env.ErrorDispatch(err)
	}
	it.SetId(newId)

	return nil
}
