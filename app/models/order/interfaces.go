// Package order represents abstraction of business layer purchase order object
package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// Package global constants
const (
	MODEL_NAME_ORDER            = "Order"
	MODEL_NAME_ORDER_COLLECTION = "OrderCollection"

	MODEL_NAME_ORDER_ITEM_COLLECTION = "OrderItemCollection"
)

// I_OrderItem represents interface to access business layer implementation of purchase order item object
type I_OrderItem interface {
	GetId() string
	SetId(newId string) error

	GetName() string
	GetSku() string

	GetQty() int

	GetPrice() float64

	GetWeight() float64

	GetOptions() map[string]interface{}

	models.I_Object
}

// I_Order represents interface to access business layer implementation of purchase order object
type I_Order interface {
	GetItems() []I_OrderItem

	AddItem(productId string, qty int, productOptions map[string]interface{}) (I_OrderItem, error)
	RemoveItem(itemIdx int) error

	CalculateTotals() error

	NewIncrementId() error

	GetIncrementId() string
	SetIncrementId(incrementId string) error

	GetSubtotal() float64
	GetGrandTotal() float64

	GetDiscountAmount() float64
	GetTaxAmount() float64
	GetShippingAmount() float64

	GetShippingAddress() visitor.I_VisitorAddress
	GetBillingAddress() visitor.I_VisitorAddress

	GetShippingMethod() string
	GetPaymentMethod() string

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
}

// I_OrderCollection represents interface to access business layer implementation of purchase order collection
type I_OrderCollection interface {
	ListOrders() []I_Order

	models.I_Collection
}

// I_OrderItemCollection represents interface to access business layer implementation of purchase order item collection
type I_OrderItemCollection interface {
	models.I_Collection
}
