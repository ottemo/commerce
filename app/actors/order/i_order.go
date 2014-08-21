package order

import (
	"fmt"
	"sort"
	"strconv"

	"errors"

	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models/order"
)

// returns order items for current order
func (it *DefaultOrder) GetItems() []order.I_OrderItem {
	result := make([]order.I_OrderItem, 0)

	keys := make([]int, 0)
	for key, _ := range it.Items {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		result = append(result, it.Items[key])
	}

	return result

}

// adds line item to current order, or returns error
func (it *DefaultOrder) AddItem(productId string, qty int, productOptions map[string]interface{}) (order.I_OrderItem, error) {

	orderItem := new(DefaultOrderItem)

	orderItem.order = it

	it.maxIdx += 1
	orderItem.idx = it.maxIdx

	err := orderItem.Set("product_id", productId)
	if err != nil {
		return nil, err
	}

	orderItem.Set("qty", qty)
	if err != nil {
		return nil, err
	}

	orderItem.Set("options", productOptions)
	if err != nil {
		return nil, err
	}

	it.Items[orderItem.idx] = orderItem

	return orderItem, nil
}

// removes line item from current order, or returns error
func (it *DefaultOrder) RemoveItem(itemIdx int) error {
	if orderItem, present := it.Items[itemIdx]; present {

		dbEngine := db.GetDBEngine()
		if dbEngine == nil {
			return errors.New("can't get DB engine")
		}

		orderItemsCollection, err := dbEngine.GetCollection(ORDER_ITEMS_COLLECTION_NAME)
		if err != nil {
			return err
		}

		err = orderItemsCollection.DeleteById(orderItem.GetId())
		if err != nil {
			return err
		}

		delete(it.Items, itemIdx)

		return nil
	} else {
		return errors.New("can't find index " + strconv.Itoa(itemIdx))
	}
}

// assigns new unique increment id to order
func (it *DefaultOrder) NewIncrementId() {
	lastIncrementIdMutex.Lock()

	lastIncrementId += 1
	it.IncrementId = fmt.Sprintf(INCREMENT_ID_FORMAT, lastIncrementId)

	lastIncrementIdMutex.Unlock()
}

// recalculates order Subtotal and GrandTotal
func (it *DefaultOrder) CalculateTotals() error {

	var subtotal float64 = 0.0
	for _, orderItem := range it.Items {
		subtotal += orderItem.GetPrice() * float64(orderItem.GetQty())
	}
	it.Subtotal = subtotal

	it.GrandTotal = it.Subtotal + it.ShippingAmount + it.TaxAmount - it.Discount

	return nil
}

// returns subtotal of order
func (it *DefaultOrder) GetSubtotal() float64 {
	return it.Subtotal
}

// returns grand total of order
func (it *DefaultOrder) GetGrandTotal() float64 {
	return it.GrandTotal
}

// returns discount amount applied to order
func (it *DefaultOrder) GetDiscountAmount() float64 {
	return it.Discount
}

// returns tax amount applied to order
func (it *DefaultOrder) GetTaxAmount() float64 {
	return it.TaxAmount
}

// returns order shipping cost
func (it *DefaultOrder) GetShippingAmount() float64 {
	return it.ShippingAmount
}

// returns shipping method for order
func (it *DefaultOrder) GetShippingMethod() string {
	return it.ShippingMethod
}

// returns payment method used for order
func (it *DefaultOrder) GetPaymentMethod() string {
	return it.PaymentMethod
}
