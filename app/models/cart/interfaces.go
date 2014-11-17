// Package cart represents abstraction of business layer cart object
package cart

import (
	"github.com/ottemo/foundation/app/models"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

// Package global constants
const (
	CART_MODEL_NAME          = "Cart"
	SESSION_KEY_CURRENT_CART = "cart_id"
)

// I_CartItem represents interface to access business layer implementation of cart item object
type I_CartItem interface {
	GetId() string
	SetId(newId string) error

	GetIdx() int
	SetIdx(newIdx int) error

	GetProductId() string
	GetProduct() product.I_Product

	GetQty() int
	SetQty(qty int) error

	GetOptions() map[string]interface{}
	SetOption(optionName string, optionValue interface{}) error

	GetCart() I_Cart
}

// I_Cart represents interface to access business layer implementation of cart object
type I_Cart interface {
	AddItem(productId string, qty int, options map[string]interface{}) (I_CartItem, error)

	RemoveItem(itemIdx int) error

	SetQty(itemIdx int, qty int) error

	GetItems() []I_CartItem
	GetSubtotal() float64

	GetVisitorId() string
	SetVisitorId(string) error

	GetVisitor() visitor.I_Visitor

	Activate() error
	Deactivate() error

	IsActive() bool

	MakeCartForVisitor(visitorId string) error

	SetCartInfo(infoAttribute string, infoValue interface{}) error
	GetCartInfo() map[string]interface{}

	models.I_Model
	models.I_Storable
}
