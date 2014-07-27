package cart

import (
	"errors"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/product"
)


// returns id of cart item
func (it *DefaultCartItem) GetId() string {
	return it.id
}


// sets id to cart item
func (it *DefaultCartItem) SetId(newId string) error {
	it.id = newId
	return nil
}


// returns product id which cart item represents
func (it *DefaultCartItem) GetProductId() string {
	return it.ProductId
}


// returns product instance which cart item represents
func (it *DefaultCartItem) GetProduct() product.I_Product {
	return nil
}


// returns current cart item qty
func (it *DefaultCartItem) GetQty() int {
	return it.Qty
}


// sets qty for current cart item
func (it *DefaultCartItem) SetQty(qty int) error {
	if (qty > 0) {
		it.Qty = qty
	} else {
		return errors.New("qty must be greater then 0")
	}

	return nil
}


// returns all item options or nil
func (it *DefaultCartItem) GetOptions() map[string]interface{} {
	return it.Options
}


// set option to cart item
func (it *DefaultCartItem) SetOption(optionName string, optionValue interface{}) error {
	if it.Options == nil {
		it.Options = make(map[string]interface{})
	}

	it.Options[optionName] = optionValue

	return nil
}


// returns cart item belongs to
func (it *DefaultCartItem) GetCart() cart.I_Cart {
	return it.Cart
}
