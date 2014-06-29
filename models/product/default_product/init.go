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


	err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "attribute/list", ListProductAttributesRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "DELETE", "attribute/remove/:attribute", RemoveProductAttributeRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "POST", "attribute/add", AddProductAttributeRestAPI )
	if err != nil { return err }


	err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "media/list/:productId/:mediaType", MediaListRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "POST", "media/add/:productId/:mediaType/:mediaName", MediaAddRestAPI )
	if err != nil { return err }
	err = rest_service.GetRestService().RegisterJsonAPI("product", "DELETE", "media/remove/:productId/:mediaType/:mediaName", MediaRemoveRestAPI )
	if err != nil { return err }

	// err = rest_service.GetRestService().RegisterJsonAPI("product", "GET", "media/get/:productId/:mediaType/:mediaName", nil ) // TODO: it is not a JSON API
	// if err != nil { return err }


	return nil
}
