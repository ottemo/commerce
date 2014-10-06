package tax

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// module entry point before app start
func init() {
	instance := new(DefaultTax)

	checkout.RegisterTax(instance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			collection.AddColumn("code", "text", true)
			collection.AddColumn("country", "text", true)
			collection.AddColumn("state", "text", true)
			collection.AddColumn("zip", "text", false)
			collection.AddColumn("rate", "text", false)
		} else {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorNew("Can't get database engine")
	}

	return nil
}
