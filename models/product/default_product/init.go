package default_product

import (
	"errors"
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/database"

	"github.com/ottemo/foundation/api"
)

func init(){
	instance := new(DefaultProductModel)
	
	models.RegisterModel("Product", instance )
	database.RegisterOnDatabaseStart( instance.setupModel )

	api.RegisterOnRestServiceStart( instance.setupAPI )
}


func (it *DefaultProductModel) setupModel() error {

	if dbEngine := database.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Product"); err == nil {
			collection.AddColumn("sku", "text", true)
			collection.AddColumn("name", "text", true)
			collection.AddColumn("description", "text", true)
			collection.AddColumn("default_image", "text", true)
			collection.AddColumn("price", "numeric", true)
		}else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}


func (it *DefaultProductModel) setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("product", "GET", "list", it.ListProductsRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "GET", "get/:id", it.GetProductRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "POST", "create", it.CreateProductRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "PUT", "update/:id", it.UpdateProductRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "DELETE", "delete/:id", it.DeleteProductRestAPI )
	if err != nil { return err }


	err = api.GetRestService().RegisterAPI("product", "GET", "attribute/list", it.ListProductAttributesRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "DELETE", "attribute/remove/:attribute", it.RemoveProductAttributeRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "POST", "attribute/add", it.AddProductAttributeRestAPI )
	if err != nil { return err }


	err = api.GetRestService().RegisterAPI("product", "GET", "media/get/:productId/:mediaType/:mediaName", it.MediaGetRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "GET", "media/list/:productId/:mediaType", it.MediaListRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "POST", "media/add/:productId/:mediaType/:mediaName", it.MediaAddRestAPI )
	if err != nil { return err }
	err = api.GetRestService().RegisterAPI("product", "DELETE", "media/remove/:productId/:mediaType/:mediaName", it.MediaRemoveRestAPI )
	if err != nil { return err }


	return nil
}
