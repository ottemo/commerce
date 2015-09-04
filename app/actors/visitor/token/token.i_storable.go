package token

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetID returns the Default Visitor Address as a string
func (it *DefaultVisitorCard) GetID() string {
	return it.id
}

// SetID takes a string as input and sets the ID on the Visitor Address
func (it *DefaultVisitorCard) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load will take Visitor Address ID and retrieve it from the database
func (it *DefaultVisitorCard) Load(ID string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(ConstCollectionNameVisitorToken); err == nil {

			if values, err := collection.LoadByID(ID); err == nil {
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
func (it *DefaultVisitorCard) Delete() error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(ConstCollectionNameVisitorToken); err == nil {
			err := collection.DeleteByID(it.GetID())
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
func (it *DefaultVisitorCard) Save() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(ConstCollectionNameVisitorToken); err == nil {

			//if it.ZipCode== "" {
			//	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c0c6fe3d-1055-4e81-aa02-e33143594242", "Zip code for address - required")
			//}

			if it.Holder == "" && it.Token == "" && it.Number == "" &&
				it.ExpirationDate == "" && it.Type == "" && it.Payment == "" {

				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "21c10a0c-25c9-44f8-bc34-01551910f3e6", "address is blank")
			}

			if newID, err := collection.Save(it.ToHashMap()); err == nil {
				it.Set("_id", newID)
				return env.ErrorDispatch(err)
			}
			return env.ErrorDispatch(err)

		}
	}
	return nil
}
