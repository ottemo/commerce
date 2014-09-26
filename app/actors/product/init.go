package product

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// module entry point before app start
func init() {
	productInstance := new(DefaultProduct)
	var _ product.I_Product = productInstance
	models.RegisterModel(product.MODEL_NAME_PRODUCT, productInstance)

	collectionInstance := new(DefaultProductCollection)
	var _ product.I_ProductCollection = collectionInstance
	models.RegisterModel(product.MODEL_NAME_PRODUCT_COLLECTION, collectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {

	collection, err := db.GetCollection(COLLECTION_NAME_PRODUCT)
	if err != nil {
		return err
	}

	collection.AddColumn("sku", "text", true)
	collection.AddColumn("name", "text", true)
	collection.AddColumn("short_description", "text", false)
	collection.AddColumn("description", "text", false)
	collection.AddColumn("default_image", "text", false)
	collection.AddColumn("price", "numeric", false)
	collection.AddColumn("weight", "numeric", false)
	collection.AddColumn("options", "text", false)

	return nil
}
