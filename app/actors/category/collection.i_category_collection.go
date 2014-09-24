package category

import (
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *DefaultCategoryCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// returns list of category model items
func (it *DefaultCategoryCollection) ListCategories() []category.I_Category {
	result := make([]category.I_Category, 0)

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
