package cart

import (
	"github.com/ottemo/foundation/app/models/cart"
)

const (
	CART_COLLECTION_NAME       = "cart"
	CART_ITEMS_COLLECTION_NAME = "cart_items"
)

type DefaultCart struct {
	id string

	VisitorId string

	Info map[string]interface{}

	Items []cart.I_CartItem
}


type DefaultCartItem struct {
	id string

	ProductId string

	Qty int

	Options map[string]interface{}

	Cart cart.I_Cart
}
