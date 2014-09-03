package paypal

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	checkout.RegisterPaymentMethod(new(PayPalExpress))
	env.RegisterOnConfigStart(setupConfig)
}
