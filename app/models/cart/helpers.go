package cart

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetCartModel retrieves current InterfaceCart model implementation
func GetCartModel() (InterfaceCart, error) {
	model, err := models.GetModel(ConstCartModelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cartModel, ok := model.(InterfaceCart)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "369c29544a204958b929036a4bdf944d", "model "+model.GetImplementationName()+" is not 'InterfaceCart' capable")
	}

	return cartModel, nil
}

// GetCartModelAndSetID retrieves current InterfaceCart model implementation and sets its ID to some value
func GetCartModelAndSetID(cartID string) (InterfaceCart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.SetID(cartID)
	if err != nil {
		return cartModel, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// LoadCartByID loads cart data into current InterfaceCart model implementation
func LoadCartByID(cartID string) (InterfaceCart, error) {

	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.Load(cartID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// GetCartForVisitor loads cart for visitor or creates new one
func GetCartForVisitor(visitorID string) (InterfaceCart, error) {
	cartModel, err := GetCartModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = cartModel.MakeCartForVisitor(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cartModel, nil
}

// GetCurrentCart returns cart for current session or creates new one
func GetCurrentCart(params *api.StructAPIHandlerParams) (InterfaceCart, error) {
	sessionCartID := params.Session.Get(ConstSessionKeyCurrentCart)

	if sessionCartID != nil && sessionCartID != "" {

		// cart id was found in session - loading cart by id
		currentCart, err := LoadCartByID(utils.InterfaceToString(sessionCartID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		return currentCart, nil

	}

	// no cart id was in session, trying to get cart for visitor
	visitorID := params.Session.Get(visitor.ConstSessionKeyVisitorID)
	if visitorID != nil {
		currentCart, err := GetCartForVisitor(utils.InterfaceToString(visitorID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		params.Session.Set(ConstSessionKeyCurrentCart, currentCart.GetID())

		return currentCart, nil
	}

	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "388af60afced4f9b92e4984af79654a7", "you are not registered")
}
