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

// RegisterShippingMethod registers given shipping method in system
func RegisterShippingMethod(shippingMethod InterfaceShippingMethod) error {
	for _, registeredMethod := range registeredShippingMethods {
		if registeredMethod == shippingMethod {
			return env.ErrorNew("shipping method already registered")
		}
	}

	registeredShippingMethods = append(registeredShippingMethods, shippingMethod)

	return nil
}

// RegisterPaymentMethod registers given payment method in system
func RegisterPaymentMethod(paymentMethod InterfacePaymentMethod) error {
	for _, registeredMethod := range registeredPaymentMethods {
		if registeredMethod == paymentMethod {
			return env.ErrorNew("payment method already registered")
		}
	}

	registeredPaymentMethods = append(registeredPaymentMethods, paymentMethod)

	return nil
}

// RegisterTax registers given tax calculator in system
func RegisterTax(tax InterfaceTax) error {
	for _, registeredTax := range registeredTaxes {
		if registeredTax == tax {
			return env.ErrorNew("tax already registered")
		}
	}

	registeredTaxes = append(registeredTaxes, tax)

	return nil
}

// RegisterDiscount registers given discount calculator in system
func RegisterDiscount(discount InterfaceDiscount) error {
	for _, registeredDiscount := range registeredDiscounts {
		if registeredDiscount == discount {
			return env.ErrorNew("discount already registered")
		}
	}

	registeredDiscounts = append(registeredDiscounts, discount)

	return nil
}

// GetRegisteredShippingMethods returns list of registered shipping methods
func GetRegisteredShippingMethods() []InterfaceShippingMethod {
	return registeredShippingMethods
}

// GetRegisteredPaymentMethods returns list of registered payment methods
func GetRegisteredPaymentMethods() []InterfacePaymentMethod {
	return registeredPaymentMethods
}

// GetRegisteredTaxes returns list of registered tax calculators
func GetRegisteredTaxes() []InterfaceTax {
	return registeredTaxes
}

// GetRegisteredDiscounts returns list of registered tax calculators
func GetRegisteredDiscounts() []InterfaceDiscount {
	return registeredDiscounts
}
