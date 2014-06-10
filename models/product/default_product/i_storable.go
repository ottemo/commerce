package default_product

import ( "github.com/ottemo/foundation/database" )

func (it *DefaultProductModel) Load() error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if values, err := collection.LoadById( it.id ); err == nil {
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

func (it *DefaultProductModel) Save() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			if newId, err := collection.Save( it.ToHashMap() ); err == nil {
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
