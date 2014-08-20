package flat

import (
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/env"
)

func (it *FlatRateShipping) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_NAME))
}

func (it *FlatRateShipping) GetCode() string {
	return SHIPPING_CODE
}

func (it *FlatRateShipping) IsAllowed(checkout checkout.I_Checkout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(CONFIG_PATH_ENABLED))
}

func (it *FlatRateShipping) GetRates(checkoutObject checkout.I_Checkout) []checkout.T_ShippingRate {

	return []checkout.T_ShippingRate{
		checkout.T_ShippingRate{
			Code:  "default",
			Name:  utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_NAME)),
			Price: utils.InterfaceToFloat64(env.ConfigGetValue(CONFIG_PATH_AMOUNT)),
		}}
}
