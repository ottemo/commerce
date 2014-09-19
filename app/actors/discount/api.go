package discount

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/utils"
)

// initializes API for discount
func setupAPI() error {
	var err error = nil

	err = api.GetRestService().RegisterAPI("discount", "GET", "apply/:code", restDiscountApply)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("discount", "GET", "neglect/:code", restDiscountNeglect)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to apply discount code to current checkout
func restDiscountApply(params *api.T_APIHandlerParams) (interface{}, error) {

	couponCode := params.RequestURLParams["code"]

	if appliedCoupons, ok := params.Session.Get(SESSION_KEY_APPLIED_DISCOUNT_CODES).([]string); ok {

		if utils.IsAmongStr(couponCode, "100OFF", "50OFF") && !utils.IsInArray(couponCode, appliedCoupons) {
			appliedCoupons = append(appliedCoupons, couponCode)
			params.Session.Set(SESSION_KEY_APPLIED_DISCOUNT_CODES, appliedCoupons)
		} else {
			return nil, errors.New("coupon code already applied")
		}
	} else {
		params.Session.Set(SESSION_KEY_APPLIED_DISCOUNT_CODES, []string{couponCode})
	}

	return "ok", nil
}

// WEB REST API function to neglect(un-apply) discount code to current checkout
//   - use "*" as code to neglect all discounts
func restDiscountNeglect(params *api.T_APIHandlerParams) (interface{}, error) {

	couponCode := params.RequestURLParams["code"]

	if couponCode == "*" {
		params.Session.Set(SESSION_KEY_APPLIED_DISCOUNT_CODES, make([]string, 0))
		return "ok", nil
	}

	if appliedCoupons, ok := params.Session.Get(SESSION_KEY_APPLIED_DISCOUNT_CODES).([]string); ok {

		newAppliedCoupons := make([]string, 0)
		for _, value := range appliedCoupons {
			if value != couponCode {
				newAppliedCoupons = append(newAppliedCoupons, value)
			}
		}

		params.Session.Set(SESSION_KEY_APPLIED_DISCOUNT_CODES, newAppliedCoupons)
	}

	return "ok", nil
}
