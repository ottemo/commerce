package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	visitorInstance := new(DefaultVisitor)
	var _ visitor.InterfaceVisitor = visitorInstance
	models.RegisterModel(visitor.ConstModelNameVisitor, visitorInstance)

	visitorCollectionInstance := new(DefaultVisitorCollection)
	var _ visitor.InterfaceVisitorCollection = visitorCollectionInstance
	models.RegisterModel(visitor.ConstModelNameVisitorCollection, visitorCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("email", db.TypeWPrecision(db.ConstTypeVarchar, 150), true)
	collection.AddColumn("validate", db.TypeWPrecision(db.ConstTypeVarchar, 128), false)
	collection.AddColumn("password", db.TypeWPrecision(db.ConstTypeVarchar, 128), false)
	collection.AddColumn("first_name", db.TypeWPrecision(db.ConstTypeVarchar, 50), true)
	collection.AddColumn("last_name", db.TypeWPrecision(db.ConstTypeVarchar, 50), true)

	collection.AddColumn("facebook_id", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)
	collection.AddColumn("google_id", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)

	collection.AddColumn("billing_address_id", db.ConstTypeID, false)
	collection.AddColumn("shipping_address_id", db.ConstTypeID, false)

	collection.AddColumn("is_admin", db.ConstTypeBoolean, false)
	collection.AddColumn("created_at", db.ConstTypeDatetime, false)

	collection, err = db.GetCollection(ConstCollectionNameVisitorToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("visitor_id", db.ConstTypeID, true)
	collection.AddColumn("payment", db.TypeWPrecision(db.ConstTypeVarchar, 150), true)

	collection.AddColumn("type", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)
	collection.AddColumn("number", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)
	collection.AddColumn("expiration_date", db.ConstTypeDatetime, false)

	collection.AddColumn("holder", db.ConstTypeVarchar, false)

	collection.AddColumn("token", db.ConstTypeVarchar, true)
	collection.AddColumn("updated", db.ConstTypeDatetime, false)

	return nil
}
