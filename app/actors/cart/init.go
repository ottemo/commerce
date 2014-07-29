package cart

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	// "github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models/cart"
)



// module entry point before app start
func init() {
	instance := new(DefaultCart)

	ifce := interface{} (instance)
	if _, ok := ifce.(models.I_Model); !ok {
		panic("DefaultCart - I_Model interface not implemented")
	}
	if _, ok := ifce.(models.I_Storable); !ok {
		panic("DefaultCart - I_Storable interface not implemented")
	}
	if _, ok := ifce.(cart.I_Cart); !ok {
		panic("DefaultCart - I_Category interface not implemented")
	}

	models.RegisterModel("Cart", instance)
	db.RegisterOnDatabaseStart(instance.setupDB)

	// api.RegisterOnRestServiceStart(instance.setupAPI)
}



// DB preparations for current model implementation
func (it *DefaultCart) setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(CART_COLLECTION_NAME)
		if err != nil {
			return err
		}

		collection.AddColumn("visitor_id", "id", true)
		collection.AddColumn("info", "text", false)

		collection, err = dbEngine.GetCollection(CART_ITEMS_COLLECTION_NAME)
		if err != nil {
			return err
		}

		collection.AddColumn("idx", "int", false)
		collection.AddColumn("cart_id", "id", true)
		collection.AddColumn("product_id", "id", true)
		collection.AddColumn("qty", "int", false)
		collection.AddColumn("options", "text", false)

	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
