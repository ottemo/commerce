package cart

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// retrieves current InterfaceCart model implementation
func GetCartModel() (InterfaceCart, error) {
	model, err := models.GetModel(ConstCartModelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cartModel, ok := model.(InterfaceCart)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceCart' capable")
	}

	return cartModel, nil
}

// retrieves current InterfaceCart model implementation and sets its ID to some value
func GetCartModelAndSetId(cartId string) (InterfaceCart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.SetId(cartId)
	if err != nil {
		return cartModel, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// loads cart data into current InterfaceCart model implementation
func LoadCartById(cartId string) (InterfaceCart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.Load(cartId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// loads cart for visitor or creates new one
func GetCartForVisitor(visitorId string) (InterfaceCart, error) {
	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.MakeCartForVisitor(visitorId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// returns cart for current session or creates new one
func GetCurrentCart(params *api.StructAPIHandlerParams) (InterfaceCart, error) {
	sessionCartId := params.Session.Get(ConstSessionKeyCurrentCart)

	if sessionCartId != nil && sessionCartId != "" {

		// cart id was found in session - loading cart by id
		currentCart, err := LoadCartById(utils.InterfaceToString(sessionCartId))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		return currentCart, nil

	} else {

		// no cart id was in session, trying to get cart for visitor
		visitorId := params.Session.Get(visitor.ConstSessionKeyVisitorID)
		if visitorId != nil {
			currentCart, err := GetCartForVisitor(utils.InterfaceToString(visitorId))
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			params.Session.Set(ConstSessionKeyCurrentCart, currentCart.GetId())

			return currentCart, nil
		} else {
			return nil, env.ErrorNew("you are not registered")
		}

	}
}
