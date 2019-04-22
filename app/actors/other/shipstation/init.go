package shipstation

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/env"
)

func init() {
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}
