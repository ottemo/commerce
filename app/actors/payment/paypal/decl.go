package paypal

const (
	PAYMENT_CODE_REST    = "paypal_rest"
	PAYMENT_CODE_EXPRESS = "paypal_express"

	PAYMENT_NAME_REST    = "PayPal (REST API)"
	PAYMENT_NAME_EXPRESS = "PayPal Express"

	CONFIG_PATH_MODE = "payment.paypal.mode"
	CONFIG_PATH_USER = "payment.paypal.user"
	CONFIG_PATH_PASS = "payment.paypal.password"

	PP_EXPRESS_ENDPOINT      = "https://api-3t.sandbox.paypal.com/nvp"
	PP_EXPRESS_REDIRECT      = "https://www.sandbox.paypal.com/webscr?cmd=_express-checkout"
	PP_EXPRESS_USER          = "paypal_api1.ottemo.io"
	PP_EXPRESS_PWD           = "1407821638"
	PP_EXPRESS_SIGNATURE     = "AFcWxV21C7fd0v3bYYYRCpSSRl31AqosoWhBSaGs-CU45dQ.JdNevqah"
	PP_EXPRESS_PAYMENTACTION = "SALE"
)

type PayPalExpress struct{}

type PayPalRest struct{}
