package visitor

import (
	"github.com/ottemo/foundation/app/actors/visitor/address"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetId returns current ID of the Visitor
func (it *DefaultVisitor) GetId() string {
	return it.id
}

// SetId sets current ID of the current Visitor
func (it *DefaultVisitor) SetId(NewID string) error {
	it.id = NewID
	return nil
}

// Load will retrieve the Visitor information from database
func (it *DefaultVisitor) Load(ID string) error {

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	values, err := collection.LoadById(ID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(values)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current Visitor from the database
func (it *DefaultVisitor) Delete() error {

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	addressCollection, err := db.GetCollection(address.COLLECTION_NAME_VISITOR_ADDRESS)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	addressCollection.AddFilter("visitor_id", "=", it.GetId())
	if _, err := addressCollection.Delete(); err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteById(it.GetId())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current Visitor to the database
func (it *DefaultVisitor) Save() error {

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if it.GetId() == "" {
		collection.AddFilter("email", "=", it.GetEmail())
		n, err := collection.Count()
		if err != nil {
			return env.ErrorDispatch(err)
		}
		if n > 0 {
			return env.ErrorNew("email already exists")
		}
	}

	storableValues := it.ToHashMap()

	delete(storableValues, "billing_address")
	delete(storableValues, "shipping_address")

	/*if it.Password == "" {
		return env.ErrorNew("password can't be blank")
	}*/

	storableValues["facebook_id"] = it.FacebookId
	storableValues["google_id"] = it.GoogleId
	storableValues["password"] = it.Password
	storableValues["validate"] = it.ValidateKey

	// shipping address save
	if it.ShippingAddress != nil {
		err := it.ShippingAddress.Save()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		storableValues["shipping_address_id"] = it.ShippingAddress.GetId()
	}

	// billing address save
	if it.BillingAddress != nil {
		err := it.BillingAddress.Save()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		storableValues["billing_address_id"] = it.BillingAddress.GetId()
	}

	// saving visitor
	if newID, err := collection.Save(storableValues); err == nil {
		it.Set("_id", newID)
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
