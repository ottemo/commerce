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

			storableValues := it.ToHashMap()

			delete(storableValues, "billing_address")
			delete(storableValues, "shipping_address")

			if shippingAddress := it.GetShippingAddress(); shippingAddress != nil {
				err := shippingAddress.Save()
				if err != nil { return err }

				storableValues["shipping_address"] = shippingAddress.GetId()
			}

			if billingAddress := it.GetShippingAddress(); billingAddress != nil {
				err := billingAddress.Save()
				if err != nil { return err }

				storableValues["billing_address"] = billingAddress.GetId()
			}

			if newId, err := collection.Save( storableValues ); err == nil {
				it.Set("_id", newId)
			} else {
				return err
			}

		} else {
			return err
		}
	}
	return nil
}
