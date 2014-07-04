package default_product

import (
	"errors"
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/database"

	"github.com/ottemo/foundation/api"
)

func init(){
	instance := new(DefaultProductModel)
	
	models.RegisterModel("Product", instance )
	database.RegisterOnDatabaseStart( instance.setupModel )

	api.RegisterOnRestServiceStart( instance.setupAPI )
}


func (it *DefaultProductModel) setupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			collection.AddColumn("sku", "text", true)
			collection.AddColumn("name", "text", true)
			collection.AddColumn("description", "text", true)
			collection.AddColumn("default_image", "text", true)
			collection.AddColumn("price", "numeric", true)
		}else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
