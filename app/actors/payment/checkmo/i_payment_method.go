package checkmo

import (
	"github.com/ottemo/foundation/app/models/checkout"
)

func (it *CheckMoneyOrder) GetName() string {
	return PAYMENT_NAME
}

func (it *CheckMoneyOrder) GetCode() string {
	return PAYMENT_CODE
}

func (it *CheckMoneyOrder) IsAllowed(checkoutInstance checkout.I_Checkout) bool {
	return true
}

func (it *CheckMoneyOrder) Authorize(checkoutInstance checkout.I_Checkout) error {
	return nil
}

func (it *CheckMoneyOrder) Capture(checkoutInstance checkout.I_Checkout) error {
	return nil
}

func (it *CheckMoneyOrder) Refund(checkoutInstance checkout.I_Checkout) error {
	return nil
}

func (it *CheckMoneyOrder) Void(checkoutInstance checkout.I_Checkout) error {
	return nil
}
