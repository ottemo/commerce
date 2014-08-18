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

func (it *CheckMoneyOrder) IsAllowed(checkout checkout.I_Checkout) bool {
	return true
}

func (it *CheckMoneyOrder) Authorize() error {
	return nil
}

func (it *CheckMoneyOrder) Capture() error {
	return nil
}

func (it *CheckMoneyOrder) Refund() error {
	return nil
}

func (it *CheckMoneyOrder) Void() error {
	return nil
}
