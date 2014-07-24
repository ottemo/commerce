package product

import (
	"github.com/ottemo/foundation/app/actors/attributes"
	"github.com/ottemo/foundation/app/models"
)

func (it *DefaultProduct) GetModelName() string {
	return "Product"
}

func (it *DefaultProduct) GetImplementationName() string {
	return "DefaultProduct"
}

func (it *DefaultProduct) New() (models.I_Model, error) {

	customAttributes, err := new(attributes.CustomAttributes).Init("product")
	if err != nil {
		return nil, err
	}

	it.listExtraAtributes = make([]string, 0)

	return &DefaultProduct{CustomAttributes: customAttributes}, nil
}
