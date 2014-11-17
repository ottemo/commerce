// Package checkout is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package checkout

import (
	"github.com/ottemo/foundation/app/models/checkout"
)

// DefaultCheckout is a default implementer of I_Checkout
type DefaultCheckout struct {
	CartId    string
	VisitorId string
	OrderId   string

	SessionId string

	ShippingAddressId string
	BillingAddressId  string

	PaymentMethodCode  string
	ShippingMethodCode string

	ShippingRate checkout.T_ShippingRate

	Taxes     []checkout.T_TaxRate
	Discounts []checkout.T_Discount

	Info map[string]interface{}

	taxesCalculateFlag     bool
	discountsCalculateFlag bool
}
