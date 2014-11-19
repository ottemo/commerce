package address

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetId returns the Default Visitor Address as a string
func (it *DefaultVisitorAddress) GetId() string {
	return it.id
}

// SetId takes a string as input and sets the ID on the Visitor Address
func (it *DefaultVisitorAddress) SetId(NewID string) error {
	it.id = NewID
	return nil
}

// Load will take Visitor Address ID and retrieve it from the database
func (it *DefaultVisitorAddress) Load(ID string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(ConstCollectionNameVisitorAddress); err == nil {

			if values, err := collection.LoadById(ID); err == nil {
				if err := it.FromHashMap(values); err != nil {
					return env.ErrorDispatch(err)
				}
			} else {
				return env.ErrorDispatch(err)
			}

		} else {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// Delete will remove the Visitor Address from the database
func (it *DefaultVisitorAddress) Delete() error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(ConstCollectionNameVisitorAddress); err == nil {
			err := collection.DeleteById(it.GetId())
			if err != nil {
				return env.ErrorDispatch(err)
			}
		} else {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// Save will persiste the Visitor Address to the database
func (it *DefaultVisitorAddress) Save() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(ConstCollectionNameVisitorAddress); err == nil {

			//if it.ZipCode== "" {
			//	return env.ErrorNew("Zip code for address - required")
			//}

			if newID, err := collection.Save(it.ToHashMap()); err == nil {
				it.Set("_id", newID)
				return env.ErrorDispatch(err)
			}
			return env.ErrorDispatch(err)

		}
	}
	return nil
}
