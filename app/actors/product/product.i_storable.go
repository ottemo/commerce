package product

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetID returns current product id
func (it *DefaultProduct) GetID() string {
	return it.id
}

// SetID sets current product id
func (it *DefaultProduct) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load loads product information from DB
func (it *DefaultProduct) Load(loadID string) error {

	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadByID(loadID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current product from DB
func (it *DefaultProduct) Delete() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current product to DB
func (it *DefaultProduct) Save() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	newID, err := collection.Save(it.ToHashMap())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.SetID(newID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
