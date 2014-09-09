package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"

	"github.com/ottemo/foundation/app/helpers/listable"
	"github.com/ottemo/foundation/app/utils"
)

// returns model name
func (it *DefaultProductCollection) GetModelName() string {
	return product.MODEL_NAME_PRODUCT_COLLECTION
}

// returns model implementation name
func (it *DefaultProductCollection) GetImplementationName() string {
	return "Default" + product.MODEL_NAME_PRODUCT_COLLECTION
}

// returns new instance of model implementation object
func (it *DefaultProductCollection) New() (models.I_Model, error) {
	helperInstance := listable.NewListableHelper(
		listable.ListableHelperDelegates{
			CollectionName: DB_COLLECTION_NAME_PRODUCT,
			ValidateExtraAttributeFunc: func(attribute string) bool {
				return utils.IsAmongStr(attribute, "sku", "name", "description", "price", "default_image")
			},
			RecordToListItemFunc: listableRecordToListItemFunc,
		})

	return &DefaultProductCollection{ListableHelper: helperInstance}, nil
}
