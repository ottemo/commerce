package rts

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/db"
)

// module entry point before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
	app.OnAppStart(initListners)
	app.OnAppStart(initSalesHistory)
}

// DB preparations for current model implementation
func initListners() error {

	env.EventRegisterListener(referrerHandler)
	env.EventRegisterListener(visitsHandler)
	env.EventRegisterListener(addToCartHandler)
	env.EventRegisterListener(reachedCheckoutHandler)
	env.EventRegisterListener(purchasedHandler)
	env.EventRegisterListener(salesHandler)

	return nil
}

// DB preparations for current model implementation
func setupDB() error {

	collection, err := db.GetCollection(COLLECTION_NAME_SALES_HISTORY)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("product_id", "id", true)
	collection.AddColumn("created_at", "datetime", false)
	collection.AddColumn("count", "int", false)

	collection, err = db.GetCollection(COLLECTION_NAME_SALES)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("product_id", "id", true)
	collection.AddColumn("count", "int", false)
	collection.AddColumn("range", "varchar(21)", false)

	return nil
}
