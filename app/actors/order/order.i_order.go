package order

import (
	"fmt"
	"sort"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

// returns order items for current order
func (it *DefaultOrder) GetItems() []order.InterfaceOrderItem {
	result := make([]order.InterfaceOrderItem, 0)

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
func (it *DefaultOrder) AddItem(productID string, qty int, productOptions map[string]interface{}) (order.InterfaceOrderItem, error) {

	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.ApplyOptions(productOptions)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	newOrderItem := new(DefaultOrderItem)
	newOrderItem.OrderID = it.GetID()

	err = newOrderItem.Set("product_id", productModel.GetID())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("qty", qty)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("options", productModel.GetOptions())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("name", productModel.GetName())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("sku", productModel.GetSku())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("short_description", productModel.GetShortDescription())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("price", productModel.GetPrice())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = newOrderItem.Set("weight", productModel.GetWeight())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	it.maxIdx += 1
	newOrderItem.idx = it.maxIdx
	it.Items[newOrderItem.idx] = newOrderItem

	return newOrderItem, nil
}

// removes line item from current order, or returns error
func (it *DefaultOrder) RemoveItem(itemIdx int) error {
	if orderItem, present := it.Items[itemIdx]; present {

		dbEngine := db.GetDBEngine()
		if dbEngine == nil {
			return env.ErrorNew("can't get DB engine")
		}

		orderItemsCollection, err := dbEngine.GetCollection(ConstCollectionNameOrderItems)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = orderItemsCollection.DeleteByID(orderItem.GetID())
		if err != nil {
			return env.ErrorDispatch(err)
		}

		delete(it.Items, itemIdx)

		return nil
	} else {
		return env.ErrorNew("can't find index " + utils.InterfaceToString(itemIdx))
	}
}

// assigns new unique increment id to order
func (it *DefaultOrder) NewIncrementID() error {
	lastIncrementIDMutex.Lock()

	lastIncrementID += 1
	it.IncrementID = fmt.Sprintf(ConstIncrementIDFormat, lastIncrementID)

	env.GetConfig().SetValue(ConstConfigPathLastIncrementID, lastIncrementID)

	lastIncrementIDMutex.Unlock()

	return nil
}

// returns increment id of order
func (it *DefaultOrder) GetIncrementID() string {
	return it.IncrementID
}

// sets increment id to order
func (it *DefaultOrder) SetIncrementID(incrementID string) error {
	it.IncrementID = incrementID

	return nil
}

// recalculates order Subtotal and GrandTotal
func (it *DefaultOrder) CalculateTotals() error {

	var subtotal float64 = 0.0
	for _, orderItem := range it.Items {
		subtotal += utils.RoundPrice(orderItem.GetPrice() * float64(orderItem.GetQty()))
	}
	it.Subtotal = utils.RoundPrice(subtotal)

	it.GrandTotal = utils.RoundPrice(it.Subtotal + it.ShippingAmount + it.TaxAmount - it.Discount)

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

// returns shipping address for order
func (it *DefaultOrder) GetShippingAddress() visitor.InterfaceVisitorAddress {
	addressModel, _ := visitor.GetVisitorAddressModel()
	addressModel.FromHashMap(it.ShippingAddress)

	return addressModel
}

// returns billing address for order
func (it *DefaultOrder) GetBillingAddress() visitor.InterfaceVisitorAddress {
	addressModel, _ := visitor.GetVisitorAddressModel()
	addressModel.FromHashMap(it.BillingAddress)

	return addressModel
}
