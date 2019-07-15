package mailchimp

import (
	"github.com/ottemo/commerce/app"
	"github.com/ottemo/commerce/env"
)

func init() {
	app.OnAppStart(appStart)
	env.RegisterOnConfigStart(setupConfig)
}

func appStart() error {
	env.EventRegisterListener("checkout.success", "mailchimp.checkoutSuccessHandler", checkoutSuccessHandler)

	return nil
}
