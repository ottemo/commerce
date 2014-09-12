package order

import (
	"fmt"
	"sort"
	"strings"

	"errors"

	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

// internal usage function to update order item fields according to options
func applyItemOptions(item *DefaultOrderItem) error {

	//making sure that product existent
	productModel, err := product.LoadProductById(item.ProductId)
	if err != nil {
		return err
	}

	// taking item specified options and product options
	productOptions := productModel.GetOptions()
	itemOptions := item.GetOptions()

	// storing start price for a case of percentage price modifier
	startPrice := item.GetPrice()

	// sorting applicable product attributes according to "order" field
	optionsApplyOrder := make([]string, 0)
	for itemOptionName, _ := range itemOptions {

		// looking only for options that customer customer set for item
		if productOption, present := productOptions[itemOptionName]; present {
			if productOption, ok := productOption.(map[string]interface{}); ok {

				orderValue := int(^uint(0) >> 1) // default order - max integer
				if optionValue, present := productOption["order"]; present {
					orderValue = utils.InterfaceToInt(optionValue)
				}

				// encoding key order to string "000000000000001 [attribute name]"
				// for future sort as string (16 digits - max for js integer)
				key := fmt.Sprintf("%.16d %s", orderValue, itemOptionName)
				optionsApplyOrder = append(optionsApplyOrder, key)
			}
		}
	}
	sort.Strings(optionsApplyOrder)

	// function to modify orderItem according to option values
	applyOptionModifiers := func(optionToApply map[string]interface{}) {
		// price modifier
		if optionValue, present := optionToApply["price"]; present {
			addPrice := utils.InterfaceToFloat64(optionValue)
			priceType := utils.InterfaceToString(optionToApply["price_type"])

			if priceType == "percent" || priceType == "%" {
				item.Price += startPrice * addPrice / 100
			} else {
				item.Price += addPrice
			}
		}

		// sku modifier
		if optionValue, present := optionToApply["sku"]; present {
			skuModifier := utils.InterfaceToString(optionValue)
			item.Sku += "-" + skuModifier
		}
	}

	// loop over item applied option in right order
	for _, itemOptionName := range optionsApplyOrder {
		// decoding key order after sort
		itemOptionName := itemOptionName[strings.Index(itemOptionName, " ")+1:]
		itemOptionValue := itemOptions[itemOptionName]

		if productOption, present := productOptions[itemOptionName]; present {
			if productOptions, ok := productOption.(map[string]interface{}); ok {

				// product option itself can contain price, sku modifiers
				applyOptionModifiers(productOptions)

				// if product option value have predefined option values, then checking their modifiers
				if productOptionValues, present := productOptions["options"]; present {
					if productOptionValues, ok := productOptionValues.(map[string]interface{}); ok {

						// option user set can be single on multi-value
						// making it uniform
						itemOptionValueSet := []string{}
						switch typedOptionValue := itemOptionValue.(type) {
						case string:
							itemOptionValueSet = []string{typedOptionValue}
						case []string:
							itemOptionValueSet = typedOptionValue
						default:
							return errors.New("invalid value for " + itemOptionName)
						}

						// loop through option values customer set for product
						for _, itemOptionValue := range itemOptionValueSet {

							if productOptionValue, present := productOptionValues[itemOptionValue]; present {
								if productOptionValue, ok := productOptionValue.(map[string]interface{}); ok {
									applyOptionModifiers(productOptionValue)
								}
							} else {
								return errors.New("invalid '" + itemOptionName + "' option value: '" + itemOptionValue)
							}

						}
					}
				}
			}
		}
	}

	return nil
}

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

	productModel, err := product.LoadProductById(productId)
	if err != nil {
		return nil, err
	}

	newOrderItem := new(DefaultOrderItem)
	newOrderItem.OrderId = it.GetId()

	err = newOrderItem.Set("product_id", productModel.GetId())
	if err != nil {
		return nil, err
	}

	err = newOrderItem.Set("qty", qty)
	if err != nil {
		return nil, err
	}

	err = newOrderItem.Set("options", productOptions)
	if err != nil {
		return nil, err
	}

	err = newOrderItem.Set("name", productModel.GetName())
	if err != nil {
		return nil, err
	}

	err = newOrderItem.Set("sku", productModel.GetSku())
	if err != nil {
		return nil, err
	}

	err = newOrderItem.Set("short_description", productModel.GetShortDescription())
	if err != nil {
		return nil, err
	}

	err = newOrderItem.Set("price", productModel.GetPrice())
	if err != nil {
		return nil, err
	}

	err = newOrderItem.Set("weight", productModel.GetWeight())
	if err != nil {
		return nil, err
	}

	err = applyItemOptions(newOrderItem)
	if err != nil {
		return nil, err
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
		return errors.New("can't find index " + utils.InterfaceToString(itemIdx))
	}
}

// assigns new unique increment id to order
func (it *DefaultOrder) NewIncrementId() error {
	lastIncrementIdMutex.Lock()

	lastIncrementId += 1
	it.IncrementId = fmt.Sprintf(INCREMENT_ID_FORMAT, lastIncrementId)

	env.GetConfig().SetValue(CONFIG_PATH_LAST_INCREMENT_ID, lastIncrementId)

	lastIncrementIdMutex.Unlock()

	return nil
}

// returns increment id of order
func (it *DefaultOrder) GetIncrementId() string {
	return it.IncrementId
}

// sets increment id to order
func (it *DefaultOrder) SetIncrementId(incrementId string) error {
	it.IncrementId = incrementId

	return nil
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

// returns shipping address for order
func (it *DefaultOrder) GetShippingAddress() visitor.I_VisitorAddress {
	addressModel, _ := visitor.GetVisitorAddressModel()
	addressModel.FromHashMap(it.ShippingAddress)

	return addressModel
}

// returns billing address for order
func (it *DefaultOrder) GetBillingAddress() visitor.I_VisitorAddress {
	addressModel, _ := visitor.GetVisitorAddressModel()
	addressModel.FromHashMap(it.BillingAddress)

	return addressModel
}
