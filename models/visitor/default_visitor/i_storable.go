package default_visitor

import ( "github.com/ottemo/foundation/database" )

func (it *DefaultVisitor) Load(Id string) error {
	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_COLLECTION_NAME ); err == nil {

			if values, err := collection.LoadById( Id ); err == nil {
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

func (it *DefaultVisitor) Save() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_COLLECTION_NAME ); err == nil {

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
