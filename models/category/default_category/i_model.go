package default_category

import(
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

func (it *DefaultCategory) GetModelName() string {
	return "Category"
}

func (it *DefaultCategory) GetImplementationName() string {
	return "DefaultCategory"
}

func (it *DefaultCategory) New() (models.I_Model, error) {
	return &DefaultCategory{ Products: make([]product.I_Product,0)  }, nil
}
