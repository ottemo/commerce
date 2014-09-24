package discount

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
)

// module entry point before app start
func init() {
	instance := new(DefaultDiscount)

	checkout.RegisterDiscount(instance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {

	collection, err := db.GetCollection(COLLECTION_NAME_COUPON_DISCOUNTS)
	if err != nil {
		collection.AddColumn("code", "text", true)
		collection.AddColumn("name", "text", false)
		collection.AddColumn("amount", "decimal", false)
		collection.AddColumn("percent", "decimal", false)
		collection.AddColumn("times", "int", false)
		collection.AddColumn("since", "datetime", false)
		collection.AddColumn("until", "datetime", false)
	}

	return nil
}
