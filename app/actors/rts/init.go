package rts

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app"
)

// module entry point before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	app.OnAppStart(initListners)
}

// DB preparations for current model implementation
func initListners() error {

	env.EventRegisterListener(referrerHandler)
	env.EventRegisterListener(visitsHandler)
	env.EventRegisterListener(addToCartHandler)
	env.EventRegisterListener(reachedCheckoutHandler)
	env.EventRegisterListener(purchasedHandler)

	return nil
}
