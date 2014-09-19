package app

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// module entry point before app start
func init() {
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}
