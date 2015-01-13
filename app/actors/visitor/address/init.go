package address

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	visitorAddressInstance := new(DefaultVisitorAddress)
	var _ visitor.InterfaceVisitorAddress = visitorAddressInstance
	models.RegisterModel(visitor.ConstModelNameVisitorAddress, visitorAddressInstance)

	visitorAddressCollectionInstance := new(DefaultVisitorAddressCollection)
	var _ visitor.InterfaceVisitorAddressCollection = visitorAddressCollectionInstance
	models.RegisterModel(visitor.ConstModelNameVisitorAddressCollection, visitorAddressCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("visitor_id", db.ConstTypeID, false)
	collection.AddColumn("first_name", db.ConstTypeVarchar, false)
	collection.AddColumn("last_name", db.ConstTypeVarchar, false)
	collection.AddColumn("company", db.ConstTypeVarchar, false)
	collection.AddColumn("address_line1", db.ConstTypeVarchar, false)
	collection.AddColumn("address_line2", db.ConstTypeVarchar, false)
	collection.AddColumn("country", db.ConstTypeVarchar, false)
	collection.AddColumn("state", db.ConstTypeVarchar, false)
	collection.AddColumn("city", db.ConstTypeVarchar, false)
	collection.AddColumn("phone", db.ConstTypeVarchar, false)
	collection.AddColumn("zip_code", db.ConstTypeVarchar, false)

	return nil
}
