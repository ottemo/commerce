package page

import (
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// returns id for cms block
func (it *DefaultCMSPage) GetId() string {
	return it.id
}

// sets id for cms block
func (it *DefaultCMSPage) SetId(newId string) error {
	it.id = newId
	return nil
}

// loads cms block information from DB
func (it *DefaultCMSPage) Load(id string) error {
	collection, err := db.GetCollection(CMS_PAGE_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbValues, err := collection.LoadById(id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.SetId(utils.InterfaceToString(dbValues["_id"]))

	it.Identifier = utils.InterfaceToString(dbValues["identifier"])
	it.URL = utils.InterfaceToString(dbValues["url"])

	it.Title = utils.InterfaceToString(dbValues["title"])
	it.Content = utils.InterfaceToString(dbValues["content"])

	it.MetaKeywords = utils.InterfaceToString(dbValues["meta_keywords"])
	it.MetaDescription = utils.InterfaceToString(dbValues["meta_description"])

	it.CreatedAt = utils.InterfaceToTime(dbValues["created_at"])
	it.UpdatedAt = utils.InterfaceToTime(dbValues["updated_at"])

	return nil
}

// removes current cms block from DB
func (it *DefaultCMSPage) Delete() error {
	collection, err := db.GetCollection(CMS_PAGE_COLLECTION_NAME)
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
func (it *DefaultCMSPage) Save() error {
	collection, err := db.GetCollection(CMS_PAGE_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// packing data before save
	currentTime := time.Now()

	storingValues := make(map[string]interface{})

	storingValues["_id"] = it.GetId()

	storingValues["url"] = it.GetURL()

	storingValues["identifier"] = it.GetIdentifier()

	storingValues["title"] = it.GetTitle()
	storingValues["content"] = it.GetContent()

	storingValues["meta_keywords"] = it.GetMetaKeywords()
	storingValues["meta_description"] = it.GetMetaDescription()

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
