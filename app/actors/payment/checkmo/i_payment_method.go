package checkmo

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
)

// returns payment method name
func (it *CheckMoneyOrder) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTitle))
}

// returns payment method code
func (it *CheckMoneyOrder) GetCode() string {
	return ConstPaymentCode
}

// returns type of payment method
func (it *CheckMoneyOrder) GetType() string {
	return checkout.ConstPaymentTypeSimple
}

// checks for method applicability
func (it *CheckMoneyOrder) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled))
}

// makes payment method authorize operation
func (it *CheckMoneyOrder) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// makes payment method capture operation
func (it *CheckMoneyOrder) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// makes payment method refund operation
func (it *CheckMoneyOrder) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// makes payment method void operation
func (it *CheckMoneyOrder) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, nil
}
