package order

import (
	"time"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns id of current order
func (it *DefaultOrder) GetId() string {
	return it.id
}

// sets id for order
func (it *DefaultOrder) SetId(NewId string) error {
	it.id = NewId
	return nil
}

// loads order information from DB
func (it *DefaultOrder) Load(Id string) error {

	// loading order
	orderCollection, err := db.GetCollection(COLLECTION_NAME_ORDER)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	values, err := orderCollection.LoadById(Id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// initializing DefaultOrder structure
	for attribute, value := range values {
		it.Set(attribute, value)
	}

	it.Items = make(map[int]order.I_OrderItem)
	it.maxIdx = 0

	// loading order items
	orderItemsCollection, err := db.GetCollection(COLLECTION_NAME_ORDER_ITEMS)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	orderItemsCollection.AddFilter("order_id", "=", it.GetId())
	orderItems, err := orderItemsCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, orderItemValues := range orderItems {
		orderItem := new(DefaultOrderItem)

		for attribute, value := range orderItemValues {
			orderItem.Set(attribute, value)
		}

		it.Items[orderItem.idx] = orderItem
	}

	return nil
}

// removes current order from DB
func (it *DefaultOrder) Delete() error {
	if it.GetId() == "" {
		return env.ErrorNew("order id is not set")
	}

	// deleting order items
	orderItemsCollection, err := db.GetCollection(COLLECTION_NAME_ORDER_ITEMS)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = orderItemsCollection.AddFilter("order_id", "=", it.GetId())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	_, err = orderItemsCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting order
	orderCollection, err := db.GetCollection(COLLECTION_NAME_ORDER)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = orderCollection.DeleteById(it.GetId())

	return env.ErrorDispatch(err)
}

// stores current order in DB
func (it *DefaultOrder) Save() error {

	orderCollection, err := db.GetCollection(COLLECTION_NAME_ORDER)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	orderItemsCollection, err := db.GetCollection(COLLECTION_NAME_ORDER_ITEMS)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// packing data before save
	orderStoringValues := it.ToHashMap()

	it.UpdatedAt = time.Now()

	newId, err := orderCollection.Save(orderStoringValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.SetId(newId)

	// storing order items
	for _, orderItem := range it.GetItems() {
		orderItem.Set("order_id", newId)
		orderItemStoringValues := orderItem.ToHashMap()

		newId, err := orderItemsCollection.Save(orderItemStoringValues)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		orderItem.SetId(newId)
	}

	return nil
}
