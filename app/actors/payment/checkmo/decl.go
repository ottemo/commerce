// Package checkmo is a "Check Money Order" implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package checkmo

// Package global constants
const (
	ConstPaymentCode = "checkmo"
	ConstPaymentName = "Check/Money Order"

	ConstConfigPathGroup   = "payment.checkmo"
	ConstConfigPathEnabled = "payment.checkmo.enabled"
	ConstConfigPathTitle   = "payment.checkmo.title"
)

// CheckMoneyOrder is a simplest implementer of InterfacePaymentMethod
type CheckMoneyOrder struct{}
