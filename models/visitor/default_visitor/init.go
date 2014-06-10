package default_visitor

import (
	"errors"
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/database"
)

func init(){
	instance := new(DefaultVisitor)

	models.RegisterModel("Visitor", instance )
	database.RegisterOnDatabaseStart( instance.SetupModel )
}


func (it *DefaultVisitor) SetupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_COLLECTION_NAME ); err == nil {
			collection.AddColumn("email", "text", true)
			collection.AddColumn("first_name", "text", true)
			collection.AddColumn("last_name", "text", true)
			collection.AddColumn("billing_address", "int", true)
			collection.AddColumn("shipping_address", "int", true)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
