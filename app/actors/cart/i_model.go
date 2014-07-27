package cart

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
)

func (it *DefaultCart) GetModelName() string {
	return "Cart"
}

func (it *DefaultCart) GetImplementationName() string {
	return "DefaultCart"
}

func (it *DefaultCart) New() (models.I_Model, error) {
	return &DefaultCart{
		Items: make([]cart.I_CartItem, 0),
		Info: make(map[string]interface{}) }, nil
}
