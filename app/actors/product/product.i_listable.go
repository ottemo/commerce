package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// returns collection of current instance type
func (it *DefaultProduct) GetCollection() models.I_Collection {
	model, _ := models.GetModel(product.MODEL_NAME_PRODUCT_COLLECTION)
	if result, ok := model.(product.I_ProductCollection); ok {
		return result
	}

	return nil
}
