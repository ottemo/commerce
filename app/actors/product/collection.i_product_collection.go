package product

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/models/product"
)

// returns database collection
func (it *DefaultProductCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// returns array of products in model instance form
func (it *DefaultProductCollection) ListProducts() []product.I_Product {
	result := make([]product.I_Product, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, dbRecordData := range dbRecords {
		productModel, err := product.GetProductModel()
		if err != nil {
			return result
		}
		productModel.FromHashMap(dbRecordData)

		result = append(result, productModel)
	}

	return result
}
