package token

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	visitorCardInstance := new(DefaultVisitorCard)
	var _ visitor.InterfaceVisitorCard = visitorCardInstance
	models.RegisterModel(visitor.ConstModelNameVisitorCard, visitorCardInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameVisitorToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("visitor_id", db.ConstTypeID, true)
	collection.AddColumn("payment", db.TypeWPrecision(db.ConstTypeVarchar, 150), true)

	collection.AddColumn("type", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)
	collection.AddColumn("number", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)
	collection.AddColumn("expiration_date", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)

	collection.AddColumn("holder", db.ConstTypeVarchar, false)

	collection.AddColumn("token", db.ConstTypeVarchar, true)
	collection.AddColumn("updated", db.ConstTypeDatetime, false)

	return nil
}
