package authorize

const (
	PAYMENT_CODE_DPM = "authorizeNetDPM"
	PAYMENT_NAME_DPM = "Authorize.Net (Direct Post)"

	DPM_ACTION_AUTHORIZE_ONLY        = "AUTH_ONLY"
	DPM_ACTION_AUTHORIZE_AND_CAPTURE = "AUTH_CAPTURE"

	CONFIG_PATH_DPM_GROUP = "payment.authorizeNetDPM"

	CONFIG_PATH_DPM_ENABLED = "payment.authorizeNetDPM.enabled"
	CONFIG_PATH_DPM_ACTION  = "payment.authorizeNetDPM.action"
	CONFIG_PATH_DPM_TITLE   = "payment.authorizeNetDPM.title"

	CONFIG_PATH_DPM_LOGIN   = "payment.authorizeNetDPM.login"
	CONFIG_PATH_DPM_KEY     = "payment.authorizeNetDPM.key"
	CONFIG_PATH_DPM_GATEWAY = "payment.authorizeNetDPM.gateway"

	CONFIG_PATH_DPM_TEST     = "payment.authorizeNetDPM.test"
	CONFIG_PATH_DPM_DEBUG    = "payment.authorizeNetDPM.debug"
	CONFIG_PATH_DPM_CHECKOUT = "payment.authorizeNetDPM.checkout"
)

type AuthorizeNetDPM struct{}
