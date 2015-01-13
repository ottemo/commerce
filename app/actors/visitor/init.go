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

	collection.AddColumn("email", db.ConstTypeVarchar, true)      // varchar(150)
	collection.AddColumn("validate", db.ConstTypeVarchar, false)  // varchar(128)
	collection.AddColumn("password", db.ConstTypeVarchar, false)  // varchar(128)
	collection.AddColumn("first_name", db.ConstTypeVarchar, true) // varchar(50)
	collection.AddColumn("last_name", db.ConstTypeVarchar, true)  // varchar(50)

	collection.AddColumn("facebook_id", db.ConstTypeVarchar, true) // varchar(100)
	collection.AddColumn("google_id", db.ConstTypeVarchar, true)   // varchar(100)

	collection.AddColumn("billing_address_id", db.ConstTypeID, false)
	collection.AddColumn("shipping_address_id", db.ConstTypeID, false)

	collection.AddColumn("is_admin", db.ConstTypeBoolean, false)
	collection.AddColumn("created_at", db.ConstTypeDatetime, false)

	return nil
}
