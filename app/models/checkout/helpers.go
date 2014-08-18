package checkout

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

// retrieves current I_Checkout model implementation
func GetCheckoutModel() (I_Checkout, error) {
	model, err := models.GetModel(CHECKOUT_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	checkoutModel, ok := model.(I_Checkout)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_Checkout' capable")
	}

	return checkoutModel, nil
}

// retrieves shipping method for given unique code or nil if no shipping method with such code
func GetShippingMethodByCode(code string) I_ShippingMehod {

	for _, shippingMethod := range ShippingMethods {
		if shippingMethod.GetCode() == code {
			return shippingMethod
		}
	}

	return nil
}

// retrieves payment method for given unique code or nil if no payment method with such code
func GetPaymentMethodByCode(code string) I_PaymentMethod {

	for _, paymentMethod := range PaymentMethods {
		if paymentMethod.GetCode() == code {
			return paymentMethod
		}
	}

	return nil
}
