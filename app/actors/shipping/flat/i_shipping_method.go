package flat

import (
	"github.com/ottemo/foundation/app/models/checkout"
)

func (it *FlatRateShipping) GetName() string {
	return SHIPPING_NAME
}

func (it *FlatRateShipping) GetCode() string {
	return SHIPPING_CODE
}

func (it *FlatRateShipping) IsAllowed(checkout checkout.I_Checkout) bool {
	return true
}

func (it *FlatRateShipping) GetRates(checkoutObject checkout.I_Checkout) []checkout.T_ShippingRate {
	rate := checkout.T_ShippingRate {
				Code:  "default",
				Name:  "Flat Rate",
				Days:  1,
				Price: 25 }

	return []checkout.T_ShippingRate{rate}
}
