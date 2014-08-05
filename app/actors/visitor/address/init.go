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
			collection.AddColumn("street", "text", false)
			collection.AddColumn("city", "text", false)
			collection.AddColumn("state", "text", false)
			collection.AddColumn("phone", "text", false)
			collection.AddColumn("zip_code", "text", false)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
