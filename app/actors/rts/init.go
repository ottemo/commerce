package rts

import (
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// module entry point before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
	app.OnAppStart(initListners)
	app.OnAppStart(initSalesHistory)

	seconds := time.Millisecond * VisitorOnlineSeconds * 1000
	ticker := time.NewTicker(seconds)
	go func() {
		for _ = range ticker.C {
			for sessionID := range OnlineSessions {
				delta := time.Now().Sub(OnlineSessions[sessionID].time)
				if float64(VisitorOnlineSeconds) < delta.Seconds() {
					DecreaseOnline(OnlineSessions[sessionID].referrerType)
					delete(OnlineSessions, sessionID)
				}
			}
		}
	}()
}

// DB preparations for current model implementation
func initListners() error {

	env.EventRegisterListener("", referrerHandler)
	env.EventRegisterListener("", visitsHandler)
	env.EventRegisterListener("", addToCartHandler)
	env.EventRegisterListener("", reachedCheckoutHandler)
	env.EventRegisterListener("", purchasedHandler)
	env.EventRegisterListener("", salesHandler)
	env.EventRegisterListener("", regVisitorAsOnlineHandler)
	env.EventRegisterListener("", visitorOnlineActionHandler)

	return nil
}

// DB preparations for current model implementation
func setupDB() error {

	collection, err := db.GetCollection(CollectionNameSalesHistory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("product_id", "id", true)
	collection.AddColumn("created_at", "datetime", false)
	collection.AddColumn("count", "int", false)

	collection, err = db.GetCollection(CollectionNameSales)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("product_id", "id", true)
	collection.AddColumn("count", "int", false)
	collection.AddColumn("range", "varchar(21)", false)

	collection, err = db.GetCollection(CollectionNameVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("day", "datetime", false)
	collection.AddColumn("visitors", "int", false)
	collection.AddColumn("cart", "int", false)
	collection.AddColumn("checkout", "int", false)
	collection.AddColumn("sales", "int", false)
	collection.AddColumn("details", "text", false)

	return nil
}
