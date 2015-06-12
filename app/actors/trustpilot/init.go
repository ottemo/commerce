package trustpilot

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(onAppStart)
	env.RegisterOnConfigStart(setupConfig)
}

// Function for every day checking for email sent to customers who order is already two week
func schedulerFunc(params map[string]interface{}) error {
	timeDay := time.Hour * 24
	currentTime := time.Now()
	twoWeeksAgo := currentTime.Add(-timeDay * 14)

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.AddFilter("created_at", ">=", twoWeeksAgo)
	dbOrderCollection.AddFilter("created_at", "<", twoWeeksAgo.Add(timeDay))

	dbRecords, err := dbOrderCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotEmailTemplate))
	for _, dbRecord := range dbRecords {
		order := utils.InterfaceToMap(dbRecord)
		customInfo := utils.InterfaceToMap(order["custom_info"])
		emailSent := utils.InterfaceToBool(customInfo[ConstOrderCustomInfoSentKey])
		orderStatus := utils.InterfaceToString(order["status"])

		if trustpilotLink, present := customInfo[ConstOrderCustomInfoLinkKey]; present && !emailSent && orderStatus != "new" {

			visitorMap := make(map[string]interface{})
			visitorEmail := utils.InterfaceToString(order["customer_email"])

			visitorMap["name"] = order["customer_name"]
			visitorMap["link"] = utils.InterfaceToString(trustpilotLink)

			emailToVisitor, err := utils.TextTemplate(emailTemplate,
				map[string]interface{}{"Emailinfo": visitorMap})

			if err != nil {
				return env.ErrorDispatch(err)
			}

			err = app.SendMail(visitorEmail, "Purchase feedback", emailToVisitor)
			if err != nil {
				return env.ErrorDispatch(err)
			}
			customInfo[ConstOrderCustomInfoSentKey] = true

			order["custom_info"] = customInfo
			_, err := dbOrderCollection.Save(order)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

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
