package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	instance := new(DefaultCheckout)
	var _ checkout.I_Checkout = instance
	models.RegisterModel(checkout.CHECKOUT_MODEL_NAME, instance)

	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}
