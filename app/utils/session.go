package utils

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/app/models/visitor"
)

// returns visitor for current session if registered or error
func GetCurrentVisitor(params *api.T_APIHandlerParams) (visitor.I_Visitor, error) {
	sessionVisitorId, ok := params.Session.Get(visitor.SESSION_KEY_VISITOR_ID).(string)
	if !ok {
		return nil, errors.New("not registered visitor")
	}

	visitorInstance, err := visitor.LoadVisitorById(sessionVisitorId)

	return visitorInstance, err
}

// returns checkout for current session or creates new one
func GetCurrentCheckout(params *api.T_APIHandlerParams) (checkout.I_Checkout, error) {
	sessionObject := params.Session.Get(checkout.SESSION_KEY_CURRENT_CHECKOUT)

	var checkoutInstance checkout.I_Checkout = nil

	// trying to get checkout object from session, otherwise creating new one
	if sessionCheckout, ok := sessionObject.(checkout.I_Checkout); ok {
		checkoutInstance = sessionCheckout

	} else {

		// making new checkout object
		newCheckoutInstance, err := checkout.GetCheckoutModel()
		if err != nil {
			return nil, err
		}

		// storing checkout object to session
		params.Session.Set(checkout.SESSION_KEY_CURRENT_CHECKOUT, newCheckoutInstance)

		//setting session
		newCheckoutInstance.SetSession(params.Session)

		checkoutInstance = newCheckoutInstance
	}

	// updating checkout object
	//-------------------------

	// setting cart
	currentCart, err := GetCurrentCart(params)
	if err != nil {
		return checkoutInstance, err
	}
	checkoutInstance.SetCart(currentCart)

	// setting visitor
	currentVisitor, err := GetCurrentVisitor(params)
	if err != nil {
		return checkoutInstance, err
	}
	checkoutInstance.SetVisitor(currentVisitor)

	return checkoutInstance, nil
}

// returns cart for current session or creates new one
func GetCurrentCart(params *api.T_APIHandlerParams) (cart.I_Cart, error) {
	sessionCartId := params.Session.Get(cart.SESSION_KEY_CURRENT_CART)

	if sessionCartId != nil && sessionCartId != "" {

		// cart id was found in session - loading cart by id
		currentCart, err := cart.LoadCartById(InterfaceToString(sessionCartId))
		if err != nil {
			return nil, err
		}

		return currentCart, nil

	} else {

		// no cart id was in session, trying to get cart for visitor
		visitorId := params.Session.Get(visitor.SESSION_KEY_VISITOR_ID)
		if visitorId != nil {
			currentCart, err := cart.GetCartForVisitor(InterfaceToString(visitorId))
			if err != nil {
				return nil, err
			}

			params.Session.Set(cart.SESSION_KEY_CURRENT_CART, currentCart.GetId())

			return currentCart, nil
		} else {
			return nil, errors.New("you are not registered")
		}

	}
}
