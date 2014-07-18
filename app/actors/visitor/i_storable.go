package visitor

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/actors/visitor/address"
)

func (it *DefaultVisitor) GetId() string {
	return it.id
}

func (it *DefaultVisitor) SetId(NewId string) error {
	it.id = NewId
	return nil
}

func (it *DefaultVisitor) Load(Id string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			values, err := collection.LoadById(Id)
			if err != nil {
				return err
			}

			err = it.FromHashMap(values)
			if err != nil {
				return err
			}

		} else {
			return err
		}
	}
	return nil
}

func (it *DefaultVisitor) Delete(Id string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			if addressCollection, err := dbEngine.GetCollection(address.VISITOR_ADDRESS_COLLECTION_NAME); err == nil {
				addressCollection.AddFilter("visitor_id", "=", it.GetId())
				if _, err := addressCollection.Delete(); err != nil {
					return err
				}
			} else {
				return err
			}

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

func (it *DefaultVisitor) Save() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {

			storableValues := it.ToHashMap()

			delete(storableValues, "billing_address")
			delete(storableValues, "shipping_address")


			storableValues["facebook_id"] = it.FacebookId
			storableValues["google_id"] = it.GoogleId
			storableValues["password"] = it.Password
			storableValues["validate"] = it.ValidateKey


			// shipping address save
			if it.ShippingAddress != nil {
				err := it.ShippingAddress.Save()
				if err != nil {
					return err
				}

				storableValues["shipping_address_id"] = it.ShippingAddress.GetId()
			}

			// billing address save
			if it.BillingAddress != nil {
				err := it.BillingAddress.Save()
				if err != nil {
					return err
				}

				storableValues["billing_address_id"] = it.BillingAddress.GetId()
			}

			// saving visitor
			if newId, err := collection.Save(storableValues); err == nil {
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
