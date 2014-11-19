package category

import (
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns database collection
func (it *DefaultCategoryCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListCategories returns list of category model items
func (it *DefaultCategoryCollection) ListCategories() []category.InterfaceCategory {
	result := make([]category.InterfaceCategory, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		categoryModel, err := category.GetCategoryModel()
		if err != nil {
			return result
		}
		categoryModel.FromHashMap(recordData)

		result = append(result, categoryModel)
	}

	return result
}
