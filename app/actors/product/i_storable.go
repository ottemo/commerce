package product

import (
	"github.com/ottemo/foundation/db"
)

func (it *DefaultProduct) GetId() string {
	return it.id
}

func (it *DefaultProduct) SetId(NewId string) error {
	it.id = NewId
	return nil
}

func (it *DefaultProduct) Load(loadId string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if values, err := collection.LoadById(loadId); err == nil {
				if err := it.FromHashMap(values); err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (it *DefaultProduct) Delete(Id string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			err := collection.DeleteById(Id)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (it *DefaultProduct) Save() error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if newId, err := collection.Save(it.ToHashMap()); err == nil {
				it.Set("_id", newId)
				return err
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
