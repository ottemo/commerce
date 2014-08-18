package checkout

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/api"
)

// sets shipping address for checkout
func (it *DefaultCheckout) SetShippingAddress(address visitor.I_VisitorAddress) error {
	it.ShippingAddressId = address.GetId()
	return nil
}

// returns checkout shipping address
func (it *DefaultCheckout) GetShippingAddress() visitor.I_VisitorAddress {
	shippingAddress, _ := visitor.LoadVisitorAddressById(it.ShippingAddressId)
	return shippingAddress
}

// sets billing address for checkout
func (it *DefaultCheckout) SetBillingAddress(address visitor.I_VisitorAddress) error {
	it.BillingAddressId = address.GetId()
	return nil
}

// returns checkout billing address
func (it *DefaultCheckout) GetBillingAddress() visitor.I_VisitorAddress {
	billingAddress, _ := visitor.LoadVisitorAddressById(it.BillingAddressId)
	return billingAddress
}

// sets payment method for checkout
func (it *DefaultCheckout) SetPaymentMethod(paymentMethod checkout.I_PaymentMethod) error {
	it.PaymentMethod = paymentMethod
	return nil
}

// returns checkout payment method
func (it *DefaultCheckout) GetPaymentMethod() checkout.I_PaymentMethod {
	return it.PaymentMethod
}

// sets payment method for checkout
func (it *DefaultCheckout) SetShippingMethod(shippingMethod checkout.I_ShippingMehod) error {
	it.ShippingMethod = shippingMethod
	return nil
}

// returns checkout shipping rate
func (it *DefaultCheckout) GetShippingRate() *checkout.T_ShippingRate {
	return it.ShippingRate
}

// sets shipping rate for checkout
func (it *DefaultCheckout) SetShippingRate(shippingRate checkout.T_ShippingRate) error {
	it.ShippingRate = &shippingRate
	return nil
}

// return checkout shipping method
func (it *DefaultCheckout) GetShippingMethod() checkout.I_ShippingMehod {
	return it.ShippingMethod
}

// sets cart for checkout
func (it *DefaultCheckout) SetCart(checkoutCart cart.I_Cart) error {
	it.CartId = checkoutCart.GetId()
	return nil
}

// return checkout cart
func (it *DefaultCheckout) GetCart() cart.I_Cart {
	cartInstance, _ := cart.LoadCartById(it.CartId)
	return cartInstance
}

// sets visitor for checkout
func (it *DefaultCheckout) SetVisitor(checkoutVisitor visitor.I_Visitor) error {
	it.VisitorId = checkoutVisitor.GetId()

	if it.BillingAddressId == "" && checkoutVisitor.GetBillingAddress() != nil {
		it.BillingAddressId = checkoutVisitor.GetBillingAddress().GetId()
	}

	if it.ShippingAddressId == "" && checkoutVisitor.GetShippingAddress() != nil {
		it.ShippingAddressId = checkoutVisitor.GetShippingAddress().GetId()
	}

	return nil
}

// return checkout visitor
func (it *DefaultCheckout) GetVisitor() visitor.I_Visitor {
	visitorInstance, _ := visitor.LoadVisitorById(it.VisitorId)
	return visitorInstance
}

// sets visitor for checkout
func (it *DefaultCheckout) SetSession(checkoutSession api.I_Session) error {
	it.Session = checkoutSession
	return nil
}

// return checkout visitor
func (it *DefaultCheckout) GetSession() api.I_Session {
	return it.Session
}

// collects taxes applied for current checkout
func (it *DefaultCheckout) GetTaxes() (float64, []checkout.T_TaxRate) {

	var amount float64 = 0

	if !it.taxesCalculateFlag {
		it.taxesCalculateFlag = true

		it.Taxes = make([]checkout.T_TaxRate, 0)
		for _, tax := range checkout.GetRegisteredTaxes() {
			for _, taxRate := range tax.CalculateTax(it) {
				it.Taxes = append(it.Taxes, taxRate)
			}
		}

		it.taxesCalculateFlag = false
	} else {
		for _, taxRate := range it.Taxes {
			amount += taxRate.Amount
		}
	}

	return amount, it.Taxes
}

// collects discounts applied for current checkout
func (it *DefaultCheckout) GetDiscounts() (float64, []checkout.T_Discount) {

	var amount float64 = 0

	if !it.discountsCalculateFlag {
		it.discountsCalculateFlag = true

		it.Discounts = make([]checkout.T_Discount, 0)
		for _, discount := range checkout.GetRegisteredDiscounts() {
			for _, discountValue := range discount.CalculateDiscount(it) {
				it.Discounts = append(it.Discounts, discountValue)
				amount += discountValue.Amount
			}
		}

		it.discountsCalculateFlag = false
	} else {
		for _, discount := range it.Discounts {
			amount += discount.Amount
		}
	}

	return amount, it.Discounts
}

// return grand total for current checkout: [cart subtotal] + [shipping rate] + [taxes] - [discounts]
func (it *DefaultCheckout) GetGrandTotal() float64 {
	var amount float64 = 0

	currentCart := it.GetCart()
	if currentCart != nil {
		amount += currentCart.GetSubtotal()
	}

	if shippingRate := it.GetShippingRate(); shippingRate != nil {
		amount += shippingRate.Price
	}

	taxAmount, _ := it.GetTaxes()
	amount += taxAmount

	discountAmount, _ := it.GetDiscounts()
	amount -= discountAmount

	return amount
}
