package cart

import (
	"errors"
	"github.com/ottemo/foundation/app/models/cart"

	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/db"
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
			return err
		}

		values, err := cartCollection.LoadById(Id)
		if err != nil {
			return err
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
			return err
		}

		cartItemsCollection.AddFilter("cart_id", "=", it.GetId())
		cartItems, err := cartItemsCollection.Load()
		if err != nil {
			return err
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
			cartItem.Options, _ = utils.DecodeJsonToStringKeyMap(cartItemValues["options"])

			it.Items[cartItem.idx] = cartItem
		}
	}

	return nil
}

// removes current cart from DB
func (it *DefaultCart) Delete() error {
	if it.GetId() == "" {
		return errors.New("cart id is not set")
	}

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("can't get DbEngine")
	}

	// deleting cart items
	cartItemsCollection, err := dbEngine.GetCollection(CART_ITEMS_COLLECTION_NAME)
	if err != nil {
		return err
	}

	err = cartItemsCollection.AddFilter("cart_id", "=", it.GetId())
	if err != nil {
		return err
	}

	_, err = cartItemsCollection.Delete()
	if err != nil {
		return err
	}

	// deleting cart
	cartCollection, err := dbEngine.GetCollection(CART_COLLECTION_NAME)
	if err != nil {
		return err
	}
	err = cartCollection.DeleteById(it.GetId())

	return err
}

// stores current cart in DB
func (it *DefaultCart) Save() error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("can't get DbEngine")
	}

	cartCollection, err := dbEngine.GetCollection(CART_COLLECTION_NAME)
	if err != nil {
		return err
	}

	cartItemsCollection, err := dbEngine.GetCollection(CART_ITEMS_COLLECTION_NAME)
	if err != nil {
		return err
	}

	// packing data before save
	cartStoringValues := make(map[string]interface{})

	cartStoringValues["_id"] = it.GetId()
	cartStoringValues["visitor_id"] = it.VisitorId
	cartStoringValues["active"] = it.Active
	cartStoringValues["info"], _ = utils.EncodeToJsonString(it.Info)

	newId, err := cartCollection.Save(cartStoringValues)
	if err != nil {
		return err
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
			return err
		}
		cartItem.SetId(newId)
	}

	return nil
}
