package flatweight

import (
	"github.com/ottemo/commerce/app/models/checkout"
	"github.com/ottemo/commerce/env"
)

func init() {
	// rates global is auto-populated via the config declaration
	env.RegisterOnConfigStart(setupConfig)

	i := new(ShippingMethod)
	if err := checkout.RegisterShippingMethod(i); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "30e9fc93-1841-4429-b5a1-c7c6cf7cd3b7", err.Error())
	}
}
