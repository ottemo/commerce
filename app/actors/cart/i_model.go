package cart

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
)



// returns model name we have implementation for
func (it *DefaultCart) GetModelName() string {
	return "Cart"
}



// returns name of current model implementation
func (it *DefaultCart) GetImplementationName() string {
	return "DefaultCart"
}



// makes new instance of model
func (it *DefaultCart) New() (models.I_Model, error) {
	return &DefaultCart{
				Items: make(map[int]cart.I_CartItem),
				Info: make(map[string]interface{}),
			}, nil
}
