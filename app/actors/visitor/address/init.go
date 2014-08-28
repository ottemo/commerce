package address

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/api"
)

// module entry point before app start
func init() {
	instance := new(DefaultVisitorAddress)

	models.RegisterModel("VisitorAddress", instance)
	db.RegisterOnDatabaseStart(instance.setupDB)

	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func (it *DefaultVisitorAddress) setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(VISITOR_ADDRESS_COLLECTION_NAME); err == nil {
			collection.AddColumn("visitor_id", "id", false)
			collection.AddColumn("first_name", "varchar(100)", false)
			collection.AddColumn("last_name", "varchar(100)", false)
			collection.AddColumn("company", "varchar(100)", false)
			collection.AddColumn("address_line1", "varchar(255)", false)
			collection.AddColumn("address_line2", "varchar(255)", false)
			collection.AddColumn("country", "varchar(50)", false)
			collection.AddColumn("state", "varchar(2)", false)
			collection.AddColumn("city", "varchar(100)", false)
			collection.AddColumn("phone", "varchar(100)", false)
			collection.AddColumn("zip_code", "varchar(10)", false)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
