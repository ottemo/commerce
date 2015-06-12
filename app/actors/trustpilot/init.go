package trustpilot

import (
	"fmt"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"time"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(onAppStart)
	env.RegisterOnConfigStart(setupConfig)
}

func schedulerFunc(params map[string]interface{}) error {
	println("hello !!!")
	fmt.Println(time.Now())
	return nil
}

// onAppStart makes module initialization on application startup
func onAppStart() error {
	if scheduler := env.GetScheduler(); scheduler != nil {
		scheduler.RegisterTask("checkOrdersToSent", schedulerFunc)
		scheduler.ScheduleRepeat("0 0 * * *", "checkOrdersToSent", nil)
	}
	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)

	return nil
}
