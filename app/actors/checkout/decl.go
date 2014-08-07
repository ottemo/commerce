package checkout

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/app/models/checkout"
)

const (
)

type DefaultCheckout struct {
	Cart    cart.I_Cart
	Visitor visitor.I_Visitor

	ShippingAddress visitor.I_VisitorAddress
	BillingAddress  visitor.I_VisitorAddress

	PaymentMethod  checkout.I_PaymentMethod
	ShippingMethod checkout.I_ShippingMehod

	Taxes     []checkout.T_TaxRate
	Discounts []checkout.T_Discount
}
