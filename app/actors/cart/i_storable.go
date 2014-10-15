package cart

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
)

// returns id of current cart
func (it *DefaultCart) GetId() string {
	return it.id
}

// sets id for cart
func (it *DefaultCart) SetId(NewId string) error {
	it.id = NewId
	return nil
}

// loads cart information from DB
func (it *DefaultCart) Load(Id string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {

		// loading category
		cartCollection, err := dbEngine.GetCollection(CART_COLLECTION_NAME)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		values, err := cartCollection.LoadById(Id)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		// initializing DefaultCart structure
		it.id = utils.InterfaceToString(values["_id"])
		it.Active = utils.InterfaceToBool(values["active"])
		it.VisitorId = utils.InterfaceToString(values["visitor_id"])
		it.Info, _ = utils.DecodeJsonToStringKeyMap(values["info"])
		it.Items = make(map[int]cart.I_CartItem)
		it.maxIdx = 0

		// loading cart items
		cartItemsCollection, err := dbEngine.GetCollection(CART_ITEMS_COLLECTION_NAME)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		cartItemsCollection.AddFilter("cart_id", "=", it.GetId())
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
			cartItem.ProductId = utils.InterfaceToString(cartItemValues["product_id"])
			cartItem.Qty = utils.InterfaceToInt(cartItemValues["qty"])
			cartItem.Options = utils.InterfaceToMap(cartItemValues["options"])

			// check for product existence and options validation
			if itemProduct := cartItem.GetProduct(); itemProduct != nil &&
				it.checkOptions(itemProduct.GetOptions(), cartItem.GetOptions()) == nil {

				it.Items[cartItem.idx] = cartItem
			} else {
				cartItemsCollection.DeleteById(utils.InterfaceToString(cartItemValues["_id"]))
			}

		}
	}

	return nil
}

// removes current cart from DB
func (it *DefaultCart) Delete() error {
	if it.GetId() == "" {
		return env.ErrorNew("cart id is not set")
	}

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("can't get DbEngine")
	}

	// deleting cart items
	cartItemsCollection, err := dbEngine.GetCollection(CART_ITEMS_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = cartItemsCollection.AddFilter("cart_id", "=", it.GetId())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	_, err = cartItemsCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting cart
	cartCollection, err := dbEngine.GetCollection(CART_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = cartCollection.DeleteById(it.GetId())

	return env.ErrorDispatch(err)
}

// stores current cart in DB
func (it *DefaultCart) Save() error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("can't get DbEngine")
	}

	cartCollection, err := dbEngine.GetCollection(CART_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	cartItemsCollection, err := dbEngine.GetCollection(CART_ITEMS_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// packing data before save
	cartStoringValues := make(map[string]interface{})

	cartStoringValues["_id"] = it.GetId()
	cartStoringValues["visitor_id"] = it.VisitorId
	cartStoringValues["active"] = it.Active
	cartStoringValues["info"] = utils.EncodeToJsonString(it.Info)

	newId, err := cartCollection.Save(cartStoringValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.SetId(newId)

	// storing cart items
	for _, cartItem := range it.GetItems() {
		cartItemStoringValues := make(map[string]interface{})

		cartItemStoringValues["_id"] = cartItem.GetId()
		cartItemStoringValues["idx"] = cartItem.GetIdx()
		cartItemStoringValues["cart_id"] = it.GetId()
		cartItemStoringValues["product_id"] = cartItem.GetProductId()
		cartItemStoringValues["qty"] = cartItem.GetQty()
		cartItemStoringValues["options"] = cartItem.GetOptions()

		newId, err := cartItemsCollection.Save(cartItemStoringValues)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		cartItem.SetId(newId)
	}

	return nil
}
