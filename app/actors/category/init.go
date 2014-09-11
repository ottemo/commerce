package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models/category"
)

// module entry point before app start
func init() {
	categoryInstance := new(DefaultCategory)
	var _ category.I_Category = categoryInstance
	models.RegisterModel(category.MODEL_NAME_CATEGORY, categoryInstance)

	categoryCollectionInstance := new(DefaultCategoryCollection)
	var _ category.I_CategoryCollection = categoryCollectionInstance
	models.RegisterModel(category.MODEL_NAME_CATEGORY_COLLECTION, categoryCollectionInstance)

	db.RegisterOnDatabaseStart(categoryInstance.setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func (it *DefaultCategory) setupDB() error {
	collection, err := db.GetCollection(COLLECTION_NAME_CATEGORY)
	if err != nil {
		return err
	}

	collection.AddColumn("parent_id", "id", true)
	collection.AddColumn("path", "text", true)
	collection.AddColumn("name", "text", true)

	collection, err = db.GetCollection(COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION)
	if err != nil {
		return err
	}

	collection.AddColumn("category_id", "id", true)
	collection.AddColumn("product_id", "id", true)

	return nil
}
