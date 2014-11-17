// Package authorize is a Authorize.Net implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package authorize

// Package global constants
const (
	// constants for transaction state
	ConstTransactionApproved      = "1"
	ConstTransactionDeclined      = "2"
	ConstTransactionError         = "3"
	ConstTransactionWaitingReview = "4"

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

// AuthorizeNetDPM is a implementer of I_PaymentMethod for a Authorize.Net Direct Post method
type AuthorizeNetDPM struct{}
