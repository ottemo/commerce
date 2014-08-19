package paypal

import (
	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	instance := new(PayPal)

	checkout.RegisterPaymentMethod(instance)

}
