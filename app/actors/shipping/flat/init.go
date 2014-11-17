package flat

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(FlatRateShipping)

	checkout.RegisterShippingMethod(instance)

	env.RegisterOnConfigStart(setupConfig)
}
