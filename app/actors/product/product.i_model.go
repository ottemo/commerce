package product

import (
	"github.com/ottemo/foundation/app/helpers/attributes"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// returns model name
func (it *DefaultProduct) GetModelName() string {
	return product.MODEL_NAME_PRODUCT
}

// returns model implementation name
func (it *DefaultProduct) GetImplementationName() string {
	return "Default" + product.MODEL_NAME_PRODUCT
}

// returns new instance of model implementation object
func (it *DefaultProduct) New() (models.I_Model, error) {

	customAttributes, err := new(attributes.CustomAttributes).Init("product")
	if err != nil {
		return nil, err
	}

	return &DefaultProduct{CustomAttributes: customAttributes}, nil
}
