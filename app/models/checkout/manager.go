package checkout

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredShippingMethods = make([]I_ShippingMethod, 0)
	registeredPaymentMethods  = make([]I_PaymentMethod, 0)

	registeredTaxes     = make([]I_Tax, 0)
	registeredDiscounts = make([]I_Discount, 0)
)

// register new shipping method to system
func RegisterShippingMethod(shippingMethod I_ShippingMethod) error {
	for _, registeredMethod := range registeredShippingMethods {
		if registeredMethod == shippingMethod {
			return env.ErrorNew("shipping method already registered")
		}
	}

	registeredShippingMethods = append(registeredShippingMethods, shippingMethod)

	return nil
}

// register new payment method to system
func RegisterPaymentMethod(paymentMethod I_PaymentMethod) error {
	for _, registeredMethod := range registeredPaymentMethods {
		if registeredMethod == paymentMethod {
			return env.ErrorNew("payment method already registered")
		}
	}

	registeredPaymentMethods = append(registeredPaymentMethods, paymentMethod)

	return nil
}

// register new tax calculator in system
func RegisterTax(tax I_Tax) error {
	for _, registeredTax := range registeredTaxes {
		if registeredTax == tax {
			return env.ErrorNew("tax already registered")
		}
	}

	registeredTaxes = append(registeredTaxes, tax)

	return nil
}

// register new discount calculator in system
func RegisterDiscount(discount I_Discount) error {
	for _, registeredDiscount := range registeredDiscounts {
		if registeredDiscount == discount {
			return env.ErrorNew("discount already registered")
		}
	}

	registeredDiscounts = append(registeredDiscounts, discount)

	return nil
}

// returns list of registered shipping methods
func GetRegisteredShippingMethods() []I_ShippingMethod {
	return registeredShippingMethods
}

// returns list of registered payment methods
func GetRegisteredPaymentMethods() []I_PaymentMethod {
	return registeredPaymentMethods
}

// returns list of registered tax calculators
func GetRegisteredTaxes() []I_Tax {
	return registeredTaxes
}

// returns list of registered tax calculators
func GetRegisteredDiscounts() []I_Discount {
	return registeredDiscounts
}
