package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

const (
	MODEL_NAME_ORDER            = "Order"
	MODEL_NAME_ORDER_COLLECTION = "OrderCollection"

	MODEL_NAME_ORDER_ITEM_COLLECTION = "OrderItemCollection"
)

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
}

type I_OrderCollection interface {
	ListOrders() []I_Order

	models.I_Collection
}

type I_OrderItemCollection interface {
	models.I_Collection
}
