// Package flatrate is a Flat Rate implementation of shipping method interface declared in
// "github.com/ottemo/commerce/app/models/checkout" package
package flatrate

import (
	"github.com/ottemo/commerce/env"
)

// Package global constants
const (
	ConstShippingCode = "flat_rate"
	ConstShippingName = "Flat Rate"

	ConstConfigPathGroup = "shipping.flat_rate"

	ConstConfigPathEnabled = "shipping.flat_rate.enabled"
	ConstConfigPathRates = "shipping.flat_rate.rates"

	ConstErrorModule = "shipping/flatrate"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// Package global variables
var (
	flatRates []interface{}
)

// ShippingMethod is a implementer of InterfaceShippingMethod for a "Flat Rate" shipping method
type ShippingMethod struct{}
