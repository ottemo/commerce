package cart

import (
	"github.com/ottemo/foundation/app/models"
)

func (it *DefaultCart) GetModelName() string {
	return "Cart"
}

func (it *DefaultCart) GetImplementationName() string {
	return "DefaultCart"
}

func (it *DefaultCart) New() (models.I_Model, error) {
	return &DefaultCart{}, nil
}
