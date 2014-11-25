// Package cart represents abstraction of business layer cart object
package cart

import (
	"github.com/ottemo/foundation/app/models"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

// Package global constants
const (
	ConstCartModelName         = "Cart"
	ConstSessionKeyCurrentCart = "cart_id"
)

// InterfaceCartItem represents interface to access business layer implementation of cart item object
type InterfaceCartItem interface {
	GetID() string
	SetID(newID string) error

	GetIdx() int
	SetIdx(newIdx int) error

	GetProductID() string
	GetProduct() product.InterfaceProduct

	GetQty() int
	SetQty(qty int) error

	GetOptions() map[string]interface{}
	SetOption(optionName string, optionValue interface{}) error

	ValidateProduct() error

	GetCart() InterfaceCart
}

// InterfaceCart represents interface to access business layer implementation of cart object
type InterfaceCart interface {
	AddItem(productID string, qty int, options map[string]interface{}) (InterfaceCartItem, error)

	RemoveItem(itemIdx int) error

	SetQty(itemIdx int, qty int) error

	GetItems() []InterfaceCartItem
	GetSubtotal() float64

	GetVisitorID() string
	SetVisitorID(string) error

	GetVisitor() visitor.InterfaceVisitor

	Activate() error
	Deactivate() error

	IsActive() bool

	MakeCartForVisitor(visitorID string) error

	SetCartInfo(infoAttribute string, infoValue interface{}) error
	GetCartInfo() map[string]interface{}

	ValidateCart() error

	models.InterfaceModel
	models.InterfaceStorable
}
