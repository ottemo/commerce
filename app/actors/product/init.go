package product

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/api"
)

func init() {
	instance := new(DefaultProduct)

	models.RegisterModel("Product", instance)
	db.RegisterOnDatabaseStart(instance.setupModel)

	api.RegisterOnRestServiceStart(instance.setupAPI)
}

func (it *DefaultProduct) setupModel() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			collection.AddColumn("sku", "text", true)
			collection.AddColumn("name", "text", true)
			collection.AddColumn("short_description", "text", false)
			collection.AddColumn("description", "text", false)
			collection.AddColumn("default_image", "text", false)
			collection.AddColumn("price", "numeric", false)
			collection.AddColumn("weight", "numeric", false)
			collection.AddColumn("size", "numeric", false)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
