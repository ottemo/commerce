package category

import (
	"github.com/ottemo/foundation/app/models/product"
)

func (it *DefaultCategory) GetName() string { return it.Name }
func (it *DefaultCategory) GetProducts() []product.I_Product {
	return it.Products
}
