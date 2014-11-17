// Package flat is a "Flat Rate" implementation of shipping method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package flat

// Package global constants
const (
	SHIPPING_CODE = "flat_rate"
	SHIPPING_NAME = "FlatRate"

	CONFIG_PATH_GROUP = "shipping.flat_rate"

	CONFIG_PATH_ENABLED = "shipping.flat_rate.enabled"
	CONFIG_PATH_AMOUNT  = "shipping.flat_rate.amount"
	CONFIG_PATH_NAME    = "shipping.flat_rate.name"
	CONFIG_PATH_DAYS    = "shipping.flat_rate.days"
)

// FlatRateShipping is a implementer of I_ShippingMethod for a "Flat Rate" shipping method
type FlatRateShipping struct{}
