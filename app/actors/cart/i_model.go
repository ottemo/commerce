package cart

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
)

// returns model name we have implementation for
func (it *DefaultCart) GetModelName() string {
	return cart.ConstCartModelName
}

// returns name of current model implementation
func (it *DefaultCart) GetImplementationName() string {
	return "DefaultCart"
}

// makes new instance of model
func (it *DefaultCart) New() (models.InterfaceModel, error) {
	return &DefaultCart{
		Items: make(map[int]cart.InterfaceCartItem),
		Info:  make(map[string]interface{}),
	}, nil
}
