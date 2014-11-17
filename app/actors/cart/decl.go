// Package cart is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/cart" package
package cart

import (
	"github.com/ottemo/foundation/app/models/cart"
)

// Package global constants
const (
	CART_COLLECTION_NAME       = "cart"
	CART_ITEMS_COLLECTION_NAME = "cart_items"
)

// DefaultCart is a default implementer of I_Cart
type DefaultCart struct {
	id string

	VisitorId string

	Info map[string]interface{}

	Items map[int]cart.I_CartItem

	Active bool

	Subtotal float64

	maxIdx int
}

// DefaultCart is a default implementer of I_Cart
type DefaultCartItem struct {
	id string

	idx int

	ProductId string

	Qty int

	Options map[string]interface{}

	Cart *DefaultCart
}
