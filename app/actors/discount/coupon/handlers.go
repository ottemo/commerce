package coupon

import (
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// checkoutSuccessHandler will add visitorID to usedCoupons by code of discount
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	orderPlaced, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4238e657-ed8e-44ed-89cf-0e66b91fecbd", "order can't be used"))
		return false
	}

	// check is discounts are applied to this order if they, make change of used coupons variable
	orderAppliedDiscounts := orderPlaced.GetDiscounts()
	visitorID := utils.InterfaceToString(orderPlaced.Get("visitor_id"))

	//	we currently not using customer email for limiting
	//	if visitorID == "" {
	//		visitorID = utils.InterfaceToString(orderPlaced.Get("customer_email"))
	//	}

	if len(orderAppliedDiscounts) > 0 && visitorID != "" {

		for _, discount := range orderAppliedDiscounts {
			if _, present := usedCoupons[discount.Code]; !present {
				usedCoupons[discount.Code] = make([]string, 0)
			}

			if !utils.IsInListStr(visitorID, usedCoupons[discount.Code]) {
				usedCoupons[discount.Code] = append(usedCoupons[discount.Code], visitorID)
			}
		}
	}

	return true
}
