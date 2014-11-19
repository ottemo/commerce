package cart

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
)

// returns id of current cart
func (it *DefaultCart) GetID() string {
	return it.id
}

// sets id for cart
func (it *DefaultCart) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// loads cart information from DB
func (it *DefaultCart) Load(ID string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {

		// loading category
		cartCollection, err := dbEngine.GetCollection(ConstCartCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		values, err := cartCollection.LoadByID(ID)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// initializing DefaultCart structure
		it.id = utils.InterfaceToString(values["_id"])
		it.Active = utils.InterfaceToBool(values["active"])
		it.VisitorID = utils.InterfaceToString(values["visitor_id"])
		it.Info, _ = utils.DecodeJSONToStringKeyMap(values["info"])
		it.Items = make(map[int]cart.InterfaceCartItem)
		it.maxIdx = 0

		// loading cart items
		cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		cartItemsCollection.AddFilter("cart_id", "=", it.GetID())
		cartItems, err := cartItemsCollection.Load()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		for _, cartItemValues := range cartItems {

			cartItem := new(DefaultCartItem)

			cartItem.id = utils.InterfaceToString(cartItemValues["_id"])
			cartItem.idx = utils.InterfaceToInt(cartItemValues["idx"])

			if cartItem.idx > it.maxIdx {
				it.maxIdx = cartItem.idx
			}

			cartItem.Cart = it
			cartItem.ProductID = utils.InterfaceToString(cartItemValues["product_id"])
			cartItem.Qty = utils.InterfaceToInt(cartItemValues["qty"])
			cartItem.Options = utils.InterfaceToMap(cartItemValues["options"])

			// check for product existence and options validation
			if itemProduct := cartItem.GetProduct(); itemProduct != nil &&
				it.checkOptions(itemProduct.GetOptions(), cartItem.GetOptions()) == nil {

				it.Items[cartItem.idx] = cartItem
			} else {
				cartItemsCollection.DeleteByID(utils.InterfaceToString(cartItemValues["_id"]))
			}

		}
	}

	return nil
}

// removes current cart from DB
func (it *DefaultCart) Delete() error {
	if it.GetID() == "" {
		return env.ErrorNew("cart id is not set")
	}

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("can't get DbEngine")
	}

	// deleting cart items
	cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = cartItemsCollection.AddFilter("cart_id", "=", it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	_, err = cartItemsCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting cart
	cartCollection, err := dbEngine.GetCollection(ConstCartCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = cartCollection.DeleteByID(it.GetID())

	return env.ErrorDispatch(err)
}

// stores current cart in DB
func (it *DefaultCart) Save() error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("can't get DbEngine")
	}

	cartCollection, err := dbEngine.GetCollection(ConstCartCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	cartItemsCollection, err := dbEngine.GetCollection(ConstCartItemsCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// packing data before save
	cartStoringValues := make(map[string]interface{})

	cartStoringValues["_id"] = it.GetID()
	cartStoringValues["visitor_id"] = it.VisitorID
	cartStoringValues["active"] = it.Active
	cartStoringValues["info"] = utils.EncodeToJSONString(it.Info)

	newID, err := cartCollection.Save(cartStoringValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.SetID(newID)

	// storing cart items
	for _, cartItem := range it.GetItems() {
		cartItemStoringValues := make(map[string]interface{})

		cartItemStoringValues["_id"] = cartItem.GetID()
		cartItemStoringValues["idx"] = cartItem.GetIdx()
		cartItemStoringValues["cart_id"] = it.GetID()
		cartItemStoringValues["product_id"] = cartItem.GetProductID()
		cartItemStoringValues["qty"] = cartItem.GetQty()
		cartItemStoringValues["options"] = cartItem.GetOptions()

		newID, err := cartItemsCollection.Save(cartItemStoringValues)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		cartItem.SetID(newID)
	}

	return nil
}
