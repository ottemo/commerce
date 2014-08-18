package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
)

const (
	ORDER_MODEL_NAME = "Order"
)

type I_OrderItem interface {
	GetId() string
	SetId(newId string) error

	GetName() string
	GetSku() string

	GetQty() int

	GetPrice() float64

	GetWeight() float64
	GetSize() float64

	GetProductOptions() map[string]interface{}

	models.I_Object
}

type I_Order interface {
	GetItems() []I_OrderItem

	AddItem(productId string, qty int, productOptions map[string]interface{}) (I_OrderItem, error)
	RemoveItem(itemIdx int) error

	CalculateTotals() error

	GetSubtotal() float64
	GetGrandTotal() float64

	GetDiscountAmount() float64
	GetTaxAmount() float64
	GetShippingAmount() float64

	GetShippingMethod() checkout.I_ShippingMehod
	GetPaymentMethod() checkout.I_PaymentMethod

	models.I_Model
	models.I_Object
	models.I_Storable
}
