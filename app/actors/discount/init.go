package discount

import (
	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	instance := new(DefaultDiscount)

	checkout.RegisterDiscount(instance)

	api.RegisterOnRestServiceStart(setupAPI)
}
