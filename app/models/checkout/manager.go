package checkout

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredShippingMethods = make([]InterfaceShippingMethod, 0)
	registeredPaymentMethods  = make([]InterfacePaymentMethod, 0)

	registeredTaxes     = make([]InterfaceTax, 0)
	registeredDiscounts = make([]InterfaceDiscount, 0)
)

// register new shipping method to system
func RegisterShippingMethod(shippingMethod InterfaceShippingMethod) error {
	for _, registeredMethod := range registeredShippingMethods {
		if registeredMethod == shippingMethod {
			return env.ErrorNew("shipping method already registered")
		}
	}

	registeredShippingMethods = append(registeredShippingMethods, shippingMethod)

	return nil
}

// register new payment method to system
func RegisterPaymentMethod(paymentMethod InterfacePaymentMethod) error {
	for _, registeredMethod := range registeredPaymentMethods {
		if registeredMethod == paymentMethod {
			return env.ErrorNew("payment method already registered")
		}
	}

	registeredPaymentMethods = append(registeredPaymentMethods, paymentMethod)

	return nil
}

// register new tax calculator in system
func RegisterTax(tax InterfaceTax) error {
	for _, registeredTax := range registeredTaxes {
		if registeredTax == tax {
			return env.ErrorNew("tax already registered")
		}
	}

	registeredTaxes = append(registeredTaxes, tax)

	return nil
}

// register new discount calculator in system
func RegisterDiscount(discount InterfaceDiscount) error {
	for _, registeredDiscount := range registeredDiscounts {
		if registeredDiscount == discount {
			return env.ErrorNew("discount already registered")
		}
	}

	registeredDiscounts = append(registeredDiscounts, discount)

	return nil
}

// returns list of registered shipping methods
func GetRegisteredShippingMethods() []InterfaceShippingMethod {
	return registeredShippingMethods
}

// returns list of registered payment methods
func GetRegisteredPaymentMethods() []InterfacePaymentMethod {
	return registeredPaymentMethods
}

// returns list of registered tax calculators
func GetRegisteredTaxes() []InterfaceTax {
	return registeredTaxes
}

// returns list of registered tax calculators
func GetRegisteredDiscounts() []InterfaceDiscount {
	return registeredDiscounts
}
