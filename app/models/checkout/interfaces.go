// Package checkout represents abstraction of business layer checkout object
package checkout

import (
	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// InterfaceCheckout represents interface to access business layer implementation of checkout object
type InterfaceCheckout interface {
	SetShippingAddress(address visitor.InterfaceVisitorAddress) error
	GetShippingAddress() visitor.InterfaceVisitorAddress

	SetBillingAddress(address visitor.InterfaceVisitorAddress) error
	GetBillingAddress() visitor.InterfaceVisitorAddress

	SetPaymentMethod(paymentMethod InterfacePaymentMethod) error
	GetPaymentMethod() InterfacePaymentMethod

	SetInfo(key string, value interface{}) error
	GetInfo(key string) interface{}

	SetShippingMethod(shippingMethod InterfaceShippingMethod) error
	GetShippingMethod() InterfaceShippingMethod

	SetShippingRate(shippingRate StructShippingRate) error
	GetShippingRate() *StructShippingRate

	GetTaxes() []StructTaxRate
	GetTaxAmount() float64

	GetDiscounts() []StructDiscount
	GetAggregatedDiscounts() []StructAggregatedDiscount
	GetDiscountAmount() float64

	GetSubtotal() float64
	GetShippingAmount() float64

	CalculateAmount(calculateTarget float64) float64
	GetGrandTotal() float64

	SetCart(checkoutCart cart.InterfaceCart) error
	GetCart() cart.InterfaceCart

	SetVisitor(checkoutVisitor visitor.InterfaceVisitor) error
	GetVisitor() visitor.InterfaceVisitor

	SetSession(api.InterfaceSession) error
	GetSession() api.InterfaceSession

	SetOrder(checkoutOrder order.InterfaceOrder) error
	GetOrder() order.InterfaceOrder

	CheckoutSuccess(checkoutOrder order.InterfaceOrder, session api.InterfaceSession) error
	SendOrderConfirmationMail() error

	Submit() (interface{}, error)

	SubmitFinish(map[string]interface{}) (interface{}, error)

	models.InterfaceModel
	models.InterfaceObject
}

// InterfaceShippingMethod represents interface to access business layer implementation of checkout shipping method
type InterfaceShippingMethod interface {
	GetName() string
	GetCode() string

	IsAllowed(checkoutInstance InterfaceCheckout) bool

	GetRates(checkoutInstance InterfaceCheckout) []StructShippingRate
}

// InterfacePaymentMethod represents interface to access business layer implementation of checkout payment method
type InterfacePaymentMethod interface {
	GetName() string
	GetCode() string
	GetType() string

	IsAllowed(checkoutInstance InterfaceCheckout) bool

	Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
	Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
	Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
	Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
}

// InterfaceTax represents interface to access business layer implementation of checkout tax system
type InterfaceTax interface {
	GetName() string
	GetCode() string

	CalculateTax(checkoutInstance InterfaceCheckout) []StructTaxRate
}

// InterfaceDiscount represents interface to access business layer implementation of checkout discount system
type InterfaceDiscount interface {
	GetName() string
	GetCode() string

	CalculateDiscount(checkoutInstance InterfaceCheckout) []StructDiscount
}

// StructShippingRate represents type to hold shipping rate information generated by implementation of InterfaceShippingMethod
type StructShippingRate struct {
	Name  string
	Code  string
	Price float64
}

// StructTaxRate represents type to hold tax rate information generated by implementation of InterfaceTax
type StructTaxRate struct {
	Name      string
	Code      string
	Amount    float64
	IsPercent bool
	Priority  float64
}

// StructDiscount represents type to hold discount information generated by implementation of InterfaceDiscount
type StructDiscount struct {
	Name      string
	Code      string
	Amount    float64
	IsPercent bool
	Priority  float64
	Object    string
	Type      string
}

// StructAggregatedDiscount represents type to hold discount information after handling in checkout calculations
type StructAggregatedDiscount struct {
	Name   string
	Code   string
	Amount float64
	Object map[string]int
	Type   string
}
