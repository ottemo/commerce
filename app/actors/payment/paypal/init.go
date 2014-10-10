package paypal

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/api"
)

// module entry point before app start
func init() {
	checkout.RegisterPaymentMethod(new(PayPalExpress))
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}
