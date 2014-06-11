package default_category

import (
	"github.com/ottemo/foundation/models/product"
)

func (it *DefaultCategory) GetName() string { return it.Name }
func (it *DefaultCategory) GetProducts() []product.I_Product {
	return it.Products
}
