package product

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	productInstance := new(DefaultProduct)
	var _ product.InterfaceProduct = productInstance
	models.RegisterModel(product.ConstModelNameProduct, productInstance)

	collectionInstance := new(DefaultProductCollection)
	var _ product.InterfaceProductCollection = collectionInstance
	models.RegisterModel(product.ConstModelNameProductCollection, collectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("sku", "text", true)
	collection.AddColumn("name", "text", true)
	collection.AddColumn("short_description", "text", false)
	collection.AddColumn("description", "text", false)
	collection.AddColumn("default_image", "text", false)
	collection.AddColumn("price", "numeric", false)
	collection.AddColumn("weight", "numeric", false)
	collection.AddColumn("options", "text", false)

	collection.AddColumn("related_pids", "[]text", false)

	return nil
}
