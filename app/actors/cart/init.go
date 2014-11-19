package cart

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models/cart"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultCart)

	ifce := interface{}(instance)
	if _, ok := ifce.(models.InterfaceModel); !ok {
		panic("DefaultCart - InterfaceModel interface not implemented")
	}
	if _, ok := ifce.(models.InterfaceStorable); !ok {
		panic("DefaultCart - InterfaceStorable interface not implemented")
	}
	if _, ok := ifce.(cart.InterfaceCart); !ok {
		panic("DefaultCart - InterfaceCategory interface not implemented")
	}

	models.RegisterModel("Cart", instance)
	db.RegisterOnDatabaseStart(instance.setupDB)

	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func (it *DefaultCart) setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(ConstCartCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("visitor_id", "id", true)
		collection.AddColumn("active", "bool", true)
		collection.AddColumn("info", "text", false)

		collection, err = dbEngine.GetCollection(ConstCartItemsCollectionName)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		collection.AddColumn("idx", "int", false)
		collection.AddColumn("cart_id", "id", true)
		collection.AddColumn("product_id", "id", true)
		collection.AddColumn("qty", "int", false)
		collection.AddColumn("options", "text", false)

	} else {
		return env.ErrorNew("Can't get database engine")
	}

	return nil
}
