package visitor

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/api"
)

// package self initializator
func init() {
	instance := new(DefaultVisitor)

	models.RegisterModel("Visitor", instance)

	db.RegisterOnDatabaseStart(instance.setupDB)

	api.RegisterOnRestServiceStart(setupAPI)
}

// setups database tables for model usage
func (it *DefaultVisitor) setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {
			collection.AddColumn("email", "id", true)
			collection.AddColumn("validate", "varchar(128)", false)
			collection.AddColumn("password", "varchar(128)", false)
			collection.AddColumn("first_name", "varchar(50)", true)
			collection.AddColumn("last_name", "varchar(50)", true)

			collection.AddColumn("facebook_id", "varchar(100)", true)
			collection.AddColumn("google_id", "varchar(100)", true)

			collection.AddColumn("billing_address_id", "id", false)
			collection.AddColumn("shipping_address_id", "id", false)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
