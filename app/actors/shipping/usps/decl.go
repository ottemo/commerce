// Package usps is a "USPS" implementation of shipping method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package usps

// Package global constants
const (
	SHIPPING_CODE = "usps"
	SHIPPING_NAME = "USPS"

	REMOVE_RATE_NAME_TAGS = true

	HTTP_ENDPOINT  = "http://production.shippingapis.com/ShippingAPI.dll"
	HTTPS_ENDPOINT = "https://secure.shippingapis.com/ShippingAPI.dll"

	CONFIG_PATH_GROUP              = "shipping.usps"
	CONFIG_PATH_ENABLED            = "shipping.usps.enabled"
	CONFIG_PATH_USER               = "shipping.usps.userid"
	CONFIG_PATH_ORIGIN_ZIP         = "shipping.usps.zip"
	CONFIG_PATH_CONTAINER          = "shipping.usps.container"
	CONFIG_PATH_SIZE               = "shipping.usps.size"
	CONFIG_PATH_DEFAULT_DIMENSIONS = "shipping.usps.default_dimensions"
	CONFIG_PATH_DEFAULT_WEIGHT     = "shipping.usps.default_weight"
	CONFIG_PATH_ALLOWED_METHODS    = "shipping.usps.allowed_methods"
	CONFIG_PATH_DEBUG_LOG          = "shipping.usps.debug_log"
)

// USPS is a implementer of I_ShippingMethod for a "USPS" shipping method
type USPS struct{}
