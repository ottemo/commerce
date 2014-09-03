package fedex

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	checkout.RegisterShippingMethod( new(FedEx) )
	env.RegisterOnConfigStart(setupConfig)
}
