package cart

import (
	"errors"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/utils"
)

// retrieves current I_Cart model implementation
func GetCartModel() (I_Cart, error) {
	model, err := models.GetModel(CART_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	cartModel, ok := model.(I_Cart)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_Cart' capable")
	}

	return cartModel, nil
}

// retrieves current I_Cart model implementation and sets its ID to some value
func GetCartModelAndSetId(cartId string) (I_Cart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, err
	}

	err = cartModel.SetId(cartId)
	if err != nil {
		return cartModel, err
	}

	return cartModel, nil
}

// loads cart data into current I_Cart model implementation
func LoadCartById(cartId string) (I_Cart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, err
	}

	err = cartModel.Load(cartId)
	if err != nil {
		return nil, err
	}

	return cartModel, nil
}

// loads cart for visitor or creates new one
func GetCartForVisitor(visitorId string) (I_Cart, error) {
	cartModel, err := GetCartModel()
	if err != nil {
		return nil, err
	}

	err = cartModel.MakeCartForVisitor(visitorId)
	if err != nil {
		return nil, err
	}

	return cartModel, nil
}

// returns cart for current session or creates new one
func GetCurrentCart(params *api.T_APIHandlerParams) (I_Cart, error) {
	sessionCartId := params.Session.Get(SESSION_KEY_CURRENT_CART)

	if sessionCartId != nil && sessionCartId != "" {

		// cart id was found in session - loading cart by id
		currentCart, err := LoadCartById(utils.InterfaceToString(sessionCartId))
		if err != nil {
			return nil, err
		}

		return currentCart, nil

	} else {

		// no cart id was in session, trying to get cart for visitor
		visitorId := params.Session.Get(visitor.SESSION_KEY_VISITOR_ID)
		if visitorId != nil {
			currentCart, err := GetCartForVisitor(utils.InterfaceToString(visitorId))
			if err != nil {
				return nil, err
			}

			params.Session.Set(SESSION_KEY_CURRENT_CART, currentCart.GetId())

			return currentCart, nil
		} else {
			return nil, errors.New("you are not registered")
		}

	}
}
