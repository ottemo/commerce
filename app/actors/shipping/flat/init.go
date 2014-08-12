package flat

import (
	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	instance := new(FlatRateShipping)

	checkout.RegisterShippingMethod(instance)

}
