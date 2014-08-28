package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
)

const ()

type DefaultCheckout struct {
	CartId    string
	VisitorId string
	OrderId   string

	Session api.I_Session

	ShippingAddressId string
	BillingAddressId  string

	PaymentMethod checkout.I_PaymentMethod

	ShippingMethod checkout.I_ShippingMehod
	ShippingRate   *checkout.T_ShippingRate

	Taxes     []checkout.T_TaxRate
	Discounts []checkout.T_Discount

	Info map[string]interface{}

	taxesCalculateFlag     bool
	discountsCalculateFlag bool
}
