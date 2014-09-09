package product

import (
	"github.com/ottemo/foundation/app/models/product"
)

func (it *DefaultProductCollection) ListProducts() []product.I_Product {
	result := make([]product.I_Product, 0)

	// loading data from DB
	collection := it.GetDBCollection()
	if collection == nil {
		return result
	}

	dbItems, err := collection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbItems {
		productModel, err := product.GetProductModel()
		if err != nil {
			return result
		}
		productModel.FromHashMap(recordData)

		result = append(result, productModel)
	}

	return result
}
