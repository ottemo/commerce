package cart

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
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
