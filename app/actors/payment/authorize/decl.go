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

	ConstPaymentCodeDPM = "authorizeNetDPM"
	ConstPaymentNameDPM = "Authorize.Net (Direct Post)"

	ConstDPMActionAuthorizeOnly       = "AUTH_ONLY"
	ConstDPMActionAuthorizeAndCapture = "AUTH_CAPTURE"

	ConstConfigPathDPMGroup = "payment.authorizeNetDPM"

	ConstConfigPathDPMEnabled = "payment.authorizeNetDPM.enabled"
	ConstConfigPathDPMAction  = "payment.authorizeNetDPM.action"
	ConstConfigPathDPMTitle   = "payment.authorizeNetDPM.title"

	ConstConfigPathDPMLogin   = "payment.authorizeNetDPM.login"
	ConstConfigPathDPMKey     = "payment.authorizeNetDPM.key"
	ConstConfigPathDPMGateway = "payment.authorizeNetDPM.gateway"

	ConstConfigPathDPMTest     = "payment.authorizeNetDPM.test"
	ConstConfigPathDPMDebug    = "payment.authorizeNetDPM.debug"
	ConstConfigPathDPMCheckout = "payment.authorizeNetDPM.checkout"
)

// AuthorizeNetDPM is a implementer of InterfacePaymentMethod for a Authorize.Net Direct Post method
type AuthorizeNetDPM struct{}
