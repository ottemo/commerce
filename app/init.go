package app

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/env"
)

// init makes package self-initialization routine
func init() {
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}
