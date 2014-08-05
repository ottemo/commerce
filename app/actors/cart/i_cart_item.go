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

// returns index value for current cart item
func (it *DefaultCartItem) GetIdx() int {
	return it.idx
}

// changes index value for current cart item if it is possible
func (it *DefaultCartItem) SetIdx(newIdx int) error {

	if newIdx < 0 {
		return errors.New("wrong cart item index")
	}

	if value, present := it.Cart.Items[newIdx]; present {
		it.Cart.Items[newIdx] = it
		it.Cart.Items[it.idx] = value
		it.idx = newIdx
	} else {
		it.Cart.Items[newIdx] = it
		it.idx = newIdx
	}

	return nil
}

// returns product id which cart item represents
func (it *DefaultCartItem) GetProductId() string {
	return it.ProductId
}

// returns product instance which cart item represents
func (it *DefaultCartItem) GetProduct() product.I_Product {
	if it.ProductId != "" {
		product, err := product.LoadProductById(it.ProductId)
		if err == nil {
			return product
		}
	}
	return nil
}

// returns current cart item qty
func (it *DefaultCartItem) GetQty() int {
	return it.Qty
}

// sets qty for current cart item
func (it *DefaultCartItem) SetQty(qty int) error {
	if qty > 0 {
		it.Qty = qty
	} else {
		return errors.New("qty must be greater then 0")
	}

	return nil
}

// removes item from the cart
func (it *DefaultCartItem) Remove() error {
	cart := it.GetCart()
	if cart != nil {
		return it.GetCart().RemoveItem(it.idx)
	} else {
		return errors.New("item is not bound to cart")
	}
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

// returns cart that item belongs to
func (it *DefaultCartItem) GetCart() cart.I_Cart {
	return it.Cart
}
