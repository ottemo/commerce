// Package flat is a Flat Rate implementation of shipping method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package flat

// Package global constants
const (
	ConstShippingCode = "flat_rate"
	ConstShippingName = "FlatRate"

	ConstConfigPathGroup = "shipping.flat_rate"

	ConstConfigPathEnabled = "shipping.flat_rate.enabled"
	ConstConfigPathAmount  = "shipping.flat_rate.amount"
	ConstConfigPathName    = "shipping.flat_rate.name"
	ConstConfigPathDays    = "shipping.flat_rate.days"
)

// FlatRateShipping is a implementer of InterfaceShippingMethod for a "Flat Rate" shipping method
type FlatRateShipping struct{}
