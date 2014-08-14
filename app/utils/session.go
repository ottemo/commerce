package utils

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/app/models/visitor"
)

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

	if checkoutInstance, ok := sessionObject.(checkout.I_Checkout); ok {

		// checkout object was found in session (no need to create new one)
		return checkoutInstance, nil

	} else {

		// making new checkout object
		newCheckoutInstance, err := checkout.GetCheckoutModel()
		if err != nil {
			return nil, err
		}

		// storing checkout object to session
		params.Session.Set(checkout.SESSION_KEY_CURRENT_CHECKOUT, newCheckoutInstance)

		// initializing created checkout object

		//setting session
		newCheckoutInstance.SetSession(params.Session)

		// setting cart
		currentCart, err := GetCurrentCart(params)
		if err != nil {
			return newCheckoutInstance, err
		}
		newCheckoutInstance.SetCart(currentCart)

		// setting visitor
		currentVisitor, err := GetCurrentVisitor(params)
		if err != nil {
			return newCheckoutInstance, err
		}
		newCheckoutInstance.SetVisitor(currentVisitor)

		return newCheckoutInstance, nil
	}
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

		if currentCheckout, err := GetCurrentCheckout(params); err == nil {
			currentCheckout.SetCart(currentCart)
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

			if currentCheckout, err := GetCurrentCheckout(params); err == nil {
				currentCheckout.SetCart(currentCart)
			}

			return currentCart, nil
		} else {
			return nil, errors.New("you are not registered")
		}

	}

	return nil, errors.New("can't get cart for current session")
}
