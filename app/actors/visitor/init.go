package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// package self initializer
func init() {
	visitorInstance := new(DefaultVisitor)
	var _ visitor.I_Visitor = visitorInstance
	models.RegisterModel(visitor.MODEL_NAME_VISITOR, visitorInstance)

	visitorCollectionInstance := new(DefaultVisitorCollection)
	var _ visitor.I_VisitorCollection = visitorCollectionInstance
	models.RegisterModel(visitor.MODEL_NAME_VISITOR_COLLECTION, visitorCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setups database tables for model usage
func setupDB() error {

	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return err
	}

	collection.AddColumn("email", "id", true)
	collection.AddColumn("validate", "varchar(128)", false)
	collection.AddColumn("password", "varchar(128)", false)
	collection.AddColumn("first_name", "varchar(50)", true)
	collection.AddColumn("last_name", "varchar(50)", true)

	collection.AddColumn("facebook_id", "varchar(100)", true)
	collection.AddColumn("google_id", "varchar(100)", true)

	collection.AddColumn("billing_address_id", "id", false)
	collection.AddColumn("shipping_address_id", "id", false)

	return nil
}
