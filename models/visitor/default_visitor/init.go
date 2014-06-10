package default_visitor

import (
	"errors"
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/database"

	"github.com/ottemo/foundation/rest_service"
)

func init(){
	instance := new(DefaultVisitor)

	models.RegisterModel("Visitor", instance )
	database.RegisterOnDatabaseStart( instance.SetupModel )

	rest_service.RegisterOnRestServiceStart( instance.SetupAPI )
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


func (it *DefaultVisitor) SetupAPI() error {
	err := rest_service.GetRestService().RegisterJsonAPI("visitor", "createVisitor", it.CreateVisitorAPI )
	if err != nil { return err }

	err = rest_service.GetRestService().RegisterJsonAPI("visitor", "updateVisitor", it.UpdateVisitorAPI )
	if err != nil { return err }

	err = rest_service.GetRestService().RegisterJsonAPI("visitor", "loadVisitor", it.LoadVisitorAPI )
	if err != nil { return err }

	return nil
}
