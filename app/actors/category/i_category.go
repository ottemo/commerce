package category

import (
	"errors"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/category"
)

func (it *DefaultCategory) GetName() string {
	return it.Name
}


func (it *DefaultCategory) GetProducts() []product.I_Product {
	return it.Products
}


func (it *DefaultCategory) GetParent() category.I_Category {
	return it.Parent
}



func (it *DefaultCategory) AddProduct(ProductId string) error {
	return errors.New("AddProduct not implemented yet")
}



func (it *DefaultCategory) RemoveProduct(ProductId string) error {
	return errors.New("RemoveProduct not implemented yet")
}
