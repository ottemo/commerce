package utils

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/app/models/visitor"
)


// returns checkout for current session or creates new one
func GetCurrentCheckout(params *api.T_APIHandlerParams) (checkout.I_Checkout, error) {
	sessionObject := params.Session.Get(checkout.SESSION_KEY_CURRENT_CHECKOUT)

	if checkoutInstance, ok := sessionObject.(checkout.I_Checkout); ok {
		return checkoutInstance, nil
	} else {
		newCheckoutInstance, err := checkout.GetCheckoutModel()
		if err != nil {
			return nil, err
		}

		params.Session.Set(checkout.SESSION_KEY_CURRENT_CHECKOUT, newCheckoutInstance)

		currentCart, err := GetCurrentCart(params)
		if err != nil {
			return newCheckoutInstance, err
		}
		newCheckoutInstance.SetCart( currentCart )

		return newCheckoutInstance, nil
	}
}



// returns cart for current session or creates new one
func GetCurrentCart(params *api.T_APIHandlerParams) (cart.I_Cart, error) {
	sessionCartId := params.Session.Get(cart.SESSION_KEY_CURRENT_CART)

	if sessionCartId != nil && sessionCartId != "" {

		currentCart, err := cart.LoadCartById(InterfaceToString(sessionCartId))
		if err != nil {
			return nil, err
		}

		return currentCart, nil
	} else {

		visitorId := params.Session.Get(visitor.SESSION_KEY_VISITOR_ID)
		if visitorId != nil {
			currentCart, err := cart.GetCartForVisitor(InterfaceToString(visitorId))
			if err != nil {
				return nil, err
			}

			params.Session.Set(cart.SESSION_KEY_CURRENT_CART, currentCart.GetId())
			params.Session.Set(checkout.SESSION_KEY_CURRENT_CHECKOUT, nil)

			return currentCart, nil
		} else {
			return nil, errors.New("you are not registered")
		}

	}

	return nil, errors.New("can't get cart for current session")
}
