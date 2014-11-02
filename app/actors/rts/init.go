package rts

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// module entry point before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
	app.OnAppStart(initListners)
	app.OnAppStart(initSalesHistory)

	seconds := time.Millisecond * VISITOR_ONLINE_SECONDS * 1000
	ticker := time.NewTicker(seconds)
	go func() {
		for _ = range ticker.C {
			for sessionId, _ := range OnlineSessions {
				delta := time.Now().Sub(OnlineSessions[sessionId].time)
				if float64(VISITOR_ONLINE_SECONDS) < delta.Seconds() {
					DecreaseOnline(OnlineSessions[sessionId].referrerType)
					delete(OnlineSessions, sessionId)
				}
			}
		}
	}()
}

// DB preparations for current model implementation
func initListners() error {

	env.EventRegisterListener(referrerHandler)
	env.EventRegisterListener(visitsHandler)
	env.EventRegisterListener(addToCartHandler)
	env.EventRegisterListener(reachedCheckoutHandler)
	env.EventRegisterListener(purchasedHandler)
	env.EventRegisterListener(salesHandler)
	env.EventRegisterListener(regVisitorAsOnlineHandler)
	env.EventRegisterListener(visitorOnlineActionHandler)

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

	collection, err = db.GetCollection(COLLECTION_NAME_VISITORS)
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
