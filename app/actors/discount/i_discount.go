package discount

import (
	"github.com/ottemo/foundation/app/models/checkout"
)

func (it *DefaultDiscount) GetName() string {
	return "CouponDiscount"
}

func (it *DefaultDiscount) GetCode() string {
	return "coupon_discount"
}

func (it *DefaultDiscount) CalculateDiscount(checkoutInstance checkout.I_Checkout) []checkout.T_Discount {

	result := make([]checkout.T_Discount, 0)

	discountableAmount := checkoutInstance.GetGrandTotal()

	if currentSession := checkoutInstance.GetSession(); currentSession != nil {

		if appliedCodes, ok := currentSession.Get(SESSION_KEY_APPLIED_DISCOUNT_CODES).([]string); ok {

			for _, discountCode := range appliedCodes {
				switch discountCode {
				case "100OFF":
					discountAmount := discountableAmount
					discountableAmount -= discountAmount

					result = append(result, checkout.T_Discount{Name: it.GetCode(), Code: discountCode, Amount: discountAmount})

				case "50OFF":
					discountAmount := discountableAmount / 2
					discountableAmount -= discountAmount

					result = append(result, checkout.T_Discount{Name: it.GetCode(), Code: discountCode, Amount: discountAmount})
				}

			}

		}
	}

	return result
}
