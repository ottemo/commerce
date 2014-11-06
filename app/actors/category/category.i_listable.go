package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
)

// returns collection of current instance type
func (it *DefaultCategory) GetCollection() models.I_Collection {
	model, _ := models.GetModel(category.MODEL_NAME_CATEGORY_COLLECTION)
	if result, ok := model.(category.I_CategoryCollection); ok {
		return result
	}

	return nil
}
