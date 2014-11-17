// Package checkmo is a "Check Money Order" implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package checkmo

// Package global constants
const (
	PAYMENT_CODE = "checkmo"
	PAYMENT_NAME = "Check/Money Order"

	CONFIG_PATH_GROUP   = "payment.checkmo"
	CONFIG_PATH_ENABLED = "payment.checkmo.enabled"
	CONFIG_PATH_TITLE   = "payment.checkmo.title"
)

// CheckMoneyOrder is a simplest implementer of I_PaymentMethod
type CheckMoneyOrder struct{}
