package usps

import (
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/env"
)

// module entry point before app start
func init() {
	instance := new(USPS)

	checkout.RegisterShippingMethod(instance)

	env.RegisterOnConfigStart(setupConfig)
}
