package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// retrieves current I_Checkout model implementation
func GetCheckoutModel() (I_Checkout, error) {
	model, err := models.GetModel(CHECKOUT_MODEL_NAME)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	checkoutModel, ok := model.(I_Checkout)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_Checkout' capable")
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

// returns checkout for current session or creates new one
func GetCurrentCheckout(params *api.T_APIHandlerParams) (I_Checkout, error) {
	sessionObject := params.Session.Get(SESSION_KEY_CURRENT_CHECKOUT)

	var checkoutInstance I_Checkout = nil

	// trying to get checkout object from session, otherwise creating new one
	if sessionCheckout, ok := sessionObject.(I_Checkout); ok {
		checkoutInstance = sessionCheckout

	} else {

		// making new checkout object
		newCheckoutInstance, err := GetCheckoutModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// storing checkout object to session
		params.Session.Set(SESSION_KEY_CURRENT_CHECKOUT, newCheckoutInstance)

		//setting session
		newCheckoutInstance.SetSession(params.Session)

		checkoutInstance = newCheckoutInstance
	}

	// updating checkout object
	//-------------------------

	// setting cart
	currentCart, err := cart.GetCurrentCart(params)
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}
	checkoutInstance.SetCart(currentCart)

	// setting visitor
	currentVisitor, err := visitor.GetCurrentVisitor(params)
	if err != nil {
		return checkoutInstance, env.ErrorDispatch(err)
	}
	checkoutInstance.SetVisitor(currentVisitor)

	return checkoutInstance, nil
}
