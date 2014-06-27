package default_product

import (
	"errors"
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/database"

	"github.com/ottemo/foundation/rest_service"
)

func init(){
	models.RegisterModel("Product", new(DefaultProductModel) )
	database.RegisterOnDatabaseStart( SetupModel )

	rest_service.RegisterOnRestServiceStart( SetupAPI )
}


func SetupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			collection.AddColumn("sku", "text", true)
			collection.AddColumn("name", "text", true)
		}else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}


func SetupAPI() error {

	var err error = nil

	err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "list", ListProductsRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "get/:id", GetProductRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "POST", "create", CreateProductRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "PUT", "update/:id", UpdateProductRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "DELETE", "delete/:id", DeleteProductRestAPI )
	if err != nil { return err }


	err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "attribute/list", nil )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "PUT", "attribute/update", nil )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "POST", "attribute/add", AddProductAttributeRestAPI )
	if err != nil { return err }


	err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "media/list", nil )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "media/get/:id/:type/:name", nil )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "POST", "media/add/:id/:type/:name", MediaAddRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "DELETE", "media/delete/:id/:type/:name", nil )
	if err != nil { return err }


	return nil
}
