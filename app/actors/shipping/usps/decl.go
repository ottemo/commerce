package usps

const (
	SHIPPING_CODE = "usps"
	SHIPPING_NAME = "USPS"

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
)

type USPS struct{}
