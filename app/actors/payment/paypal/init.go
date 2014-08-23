package paypal

import (
	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	checkout.RegisterPaymentMethod(new(PayPalRest))
	checkout.RegisterPaymentMethod(new(PayPalExpress))

}
