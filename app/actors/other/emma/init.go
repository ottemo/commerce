package emma

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/app"
	"github.com/ottemo/commerce/env"
)

func init() {
	app.OnAppStart(appStart)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}

func appStart() error {
	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)

	emmaService = *newEmmaService()

	return nil
}
