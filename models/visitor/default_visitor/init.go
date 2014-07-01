package default_visitor

import (
	"errors"
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/database"

	"github.com/ottemo/foundation/api"
)

func init(){
	instance := new(DefaultVisitor)

	models.RegisterModel("Visitor", instance )
	database.RegisterOnDatabaseStart( instance.setupModel )

	api.RegisterOnRestServiceStart( instance.setupAPI )
}


func (it *DefaultVisitor) setupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection( VISITOR_COLLECTION_NAME ); err == nil {
			collection.AddColumn("email", "text", true)
			collection.AddColumn("first_name", "text", false)
			collection.AddColumn("last_name", "text", false)
			collection.AddColumn("billing_address_id", "int", false)
			collection.AddColumn("shipping_address_id", "int", false)
		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}


func (it *DefaultVisitor) setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor", "POST", "create", it.CreateVisitorAPI )
	if err != nil { return err }

	err = api.GetRestService().RegisterAPI("visitor", "PUT", "update", it.UpdateVisitorAPI )
	if err != nil { return err }

	err = api.GetRestService().RegisterAPI("visitor", "GET", "load", it.LoadVisitorAPI )
	if err != nil { return err }

	return nil
}
