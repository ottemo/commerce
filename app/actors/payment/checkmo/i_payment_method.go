package checkmo

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
)

// returns payment method name
func (it *CheckMoneyOrder) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_TITLE))
}

// returns payment method code
func (it *CheckMoneyOrder) GetCode() string {
	return PAYMENT_CODE
}

// returns type of payment method
func (it *CheckMoneyOrder) GetType() string {
	return checkout.PAYMENT_TYPE_SIMPLE
}

// checks for method applicability
func (it *CheckMoneyOrder) IsAllowed(checkoutInstance checkout.I_Checkout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(CONFIG_PATH_ENABLED))
}

// makes payment method authorize operation
func (it *CheckMoneyOrder) Authorize(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// makes payment method capture operation
func (it *CheckMoneyOrder) Capture(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// makes payment method refund operation
func (it *CheckMoneyOrder) Refund(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// makes payment method void operation
func (it *CheckMoneyOrder) Void(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}
