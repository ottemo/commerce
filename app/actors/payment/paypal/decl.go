package paypal

const (
	PAYMENT_CODE = "paypal_express"
	PAYMENT_NAME = "PayPal Express"

	PAYMENT_ACTION_SALE = "Sale"
	PAYMENT_ACTION_AUTHORIZATION = "Authorization"

	CONFIG_PATH_GROUP = "payment.paypal"

	CONFIG_PATH_ENABLED = "payment.paypal.enabled"
	CONFIG_PATH_TITLE   = "payment.paypal.title"

	CONFIG_PATH_NVP     = "payment.paypal.nvp"
	CONFIG_PATH_GATEWAY = "payment.paypal.gateway"

	CONFIG_PATH_USER = "payment.paypal.user"
	CONFIG_PATH_PASS = "payment.paypal.password"

	CONFIG_PATH_SIGNATURE = "payment.paypal.signature"
	CONFIG_PATH_ACTION    = "payment.paypal.action"
)

type PayPalExpress struct{}

type PayPalRest struct{}
