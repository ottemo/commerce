package visitor

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/api"
)

func init() {
	instance := new(DefaultVisitor)

	models.RegisterModel("Visitor", instance)
	db.RegisterOnDatabaseStart(instance.setupModel)

	api.RegisterOnRestServiceStart(instance.setupAPI)
}

func (it *DefaultVisitor) setupModel() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_COLLECTION_NAME); err == nil {
			collection.AddColumn("email", "id", true)
			collection.AddColumn("validated", "bool", false)
			collection.AddColumn("password", "varchar(128)", false)
			collection.AddColumn("first_name", "varchar(50)", true)
			collection.AddColumn("last_name", "varchar(50)", true)
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

