// Package checkout is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package checkout

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "checkout"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultCheckout is a default implementer of InterfaceCheckout
type DefaultCheckout struct {
	CartID    string
	VisitorID string
	OrderID   string

	SessionID string

	ShippingAddress map[string]interface{}
	BillingAddress  map[string]interface{}

	PaymentMethodCode  string
	ShippingMethodCode string

	ShippingRate checkout.StructShippingRate

	Taxes     []checkout.StructTaxRate
	Discounts []checkout.StructDiscount

	Info map[string]interface{}

	// subtotal used to have static subtotal during grand total calculation
	subtotalAmount float64

	// flags enables and disables during calculation to prevent recursion
	taxesCalculateFlag      bool
	discountsCalculateFlag  bool
	grandTotalCalculateFlag bool
}
