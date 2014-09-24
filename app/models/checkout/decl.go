package checkout

const (
	CHECKOUT_MODEL_NAME = "Checkout"

	CONFIG_PATH_SHIPPING_GROUP               = "shipping"
	CONFIG_PATH_SHIPPING_ORIGIN_GROUP        = "shipping.origin"
	CONFIG_PATH_SHIPPING_ORIGIN_COUNTRY      = "shipping.origin.country"
	CONFIG_PATH_SHIPPING_ORIGIN_STATE        = "shipping.origin.state"
	CONFIG_PATH_SHIPPING_ORIGIN_CITY         = "shipping.origin.city"
	CONFIG_PATH_SHIPPING_ORIGIN_ADDRESSLINE1 = "shipping.origin.addressline1"
	CONFIG_PATH_SHIPPING_ORIGIN_ADDRESSLINE2 = "shipping.origin.addressline2"
	CONFIG_PATH_SHIPPING_ORIGIN_ZIP          = "shipping.origin.zip"

	CONFIG_PATH_PAYMENT_GROUP               = "payment"
	CONFIG_PATH_PAYMENT_ORIGIN_GROUP        = "payment.origin"
	CONFIG_PATH_PAYMENT_ORIGIN_COUNTRY      = "payment.origin.country"
	CONFIG_PATH_PAYMENT_ORIGIN_STATE        = "payment.origin.state"
	CONFIG_PATH_PAYMENT_ORIGIN_CITY         = "payment.origin.city"
	CONFIG_PATH_PAYMENT_ORIGIN_ADDRESSLINE1 = "payment.origin.addressline1"
	CONFIG_PATH_PAYMENT_ORIGIN_ADDRESSLINE2 = "payment.origin.addressline2"
	CONFIG_PATH_PAYMENT_ORIGIN_ZIP          = "payment.origin.zip"

	PAYMENT_TYPE_SIMPLE      = "simple"
	PAYMENT_TYPE_CREDIT_CARD = "cc"
	PAYMENT_TYPE_REMOTE      = "remote"
	PAYMENT_TYPE_POST        = "post"
	PAYMENT_TYPE_POST_CC     = "post_cc"

	SESSION_KEY_CURRENT_CHECKOUT = "Checkout"
)
