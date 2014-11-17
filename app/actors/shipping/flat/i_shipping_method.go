package flat

import (
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns name of shipping method
func (it *FlatRateShipping) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_NAME))
}

// GetCode returns code of shipping method
func (it *FlatRateShipping) GetCode() string {
	return SHIPPING_CODE
}

// IsAllowed checks for method applicability
func (it *FlatRateShipping) IsAllowed(checkout checkout.I_Checkout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(CONFIG_PATH_ENABLED))
}

// GetRates returns rates allowed by shipping method for a given checkout
func (it *FlatRateShipping) GetRates(checkoutObject checkout.I_Checkout) []checkout.T_ShippingRate {

	return []checkout.T_ShippingRate{
		checkout.T_ShippingRate{
			Code:  "default",
			Name:  utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_NAME)),
			Price: utils.InterfaceToFloat64(env.ConfigGetValue(CONFIG_PATH_AMOUNT)),
		}}
}
