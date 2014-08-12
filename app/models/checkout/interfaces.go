package checkout

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/app/models"
)

const (
	CHECKOUT_MODEL_NAME = "Checkout"

	SESSION_KEY_CURRENT_CHECKOUT = "Checkout"
)



type I_Checkout interface {
	SetShippingAddress(address visitor.I_VisitorAddress) error
	GetShippingAddress() visitor.I_VisitorAddress

	SetBillingAddress(address visitor.I_VisitorAddress) error
	GetBillingAddress() visitor.I_VisitorAddress

	SetPaymentMethod(paymentMethod I_PaymentMethod) error
	GetPaymentMethod() I_PaymentMethod

	SetShippingMethod(shippingMethod I_ShippingMehod) error
	GetShippingMethod() I_ShippingMehod
	SetShippingRate(shippingRate T_ShippingRate) error
	GetShippingRate() *T_ShippingRate

	SetCart(checkoutCart cart.I_Cart) error
	GetCart() cart.I_Cart

	SetVisitor(checkoutVisitor visitor.I_Visitor) error
	GetVisitor() visitor.I_Visitor

	models.I_Model
}



type I_ShippingMehod interface {
	GetName() string
	GetCode() string

	IsAllowed(checkout I_Checkout) bool

	GetRates(checkout I_Checkout) []T_ShippingRate
}



type I_PaymentMethod interface {
	GetName() string
	GetCode() string

	IsAllowed(checkout I_Checkout) bool

	Authorize() error
	Capture() error
	Refund() error
	Void()	error
}



type I_Tax interface {
	GetName() string
	GetCode() string

	CalculateTax(checkout I_Checkout) []T_TaxRate
}



type I_Discount interface {
	GetName() string
	GetCode() string

	CalculateDiscount(checkout I_Checkout) []T_Discount
}



type T_ShippingRate struct {
	Name string
	Code string
	Price float64
	Days int
}



type T_TaxRate struct {
	Name string
	Code string
	Amount float64
}



type T_Discount struct {
	Name string
	Code string
	Amount float64
}
