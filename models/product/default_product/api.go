package default_product

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

func jsonError(err error) map[string]interface{} {
	return map[string]interface{} { "error": err.Error() }
}

func AddProductAttributeRestAPI(req *http.Request, params map[string]string) map[string]interface{} {
	queryParams := req.URL.Query()

	model, err := models.GetModel("Product")
	if err != nil { return jsonError(err) }

	attribute := models.T_AttributeInfo {
		Model:      "product",
		Collection: "product",
		Attribute:  "test",
		Type:       "text",
		Label:      "Test Attribute",
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
	}


	for param, value := range queryParams {
		switch param {
		case "type":
			attribute.Type = value[0]
		case "attribute":
			attribute.Attribute = value[0]
		case "label":
			attribute.Label = value[0]
		case "group":
			attribute.Group = value[0]
		case "editors":
			attribute.Editors = value[0]
		case "options":
			attribute.Options = value[0]
		case "default":
			attribute.Default = value[0]
		}
	}


	if prod, ok := model.(models.I_CustomAttributes); ok {
		if err := prod.AddNewAttribute(attribute); err != nil {
			return jsonError( errors.New("Product new attribute error: " + err.Error()) )
		}
	} else {
		return jsonError( errors.New("product model is not I_CustomAttributes") )
	}


	return map[string]interface{} {"ok": true, "attribute": attribute}
}

// WEB REST API function used to obtain all product attributes
//   - product id must be specified in request URI "http://[site:port]/product/get/:id"
func GetProductRestAPI(req *http.Request, params map[string]string) map[string]interface{} {

	productId, isSpecifiedId := params["id"]
	if !isSpecifiedId {
		return jsonError( errors.New("product 'id' was not specified") )
	}

	if model, err := models.GetModel("Product"); err == nil {
		if model, ok := model.(product.I_Product); ok {

			err = model.Load( productId )
			if err != nil { return jsonError(err) }

			return model.ToHashMap()
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}

// WEB REST API function used to obtain product list we have in database
//   - only [_id, sku, name] attributes returns by default
func ListProductsRestAPI(req *http.Request, params map[string]string) map[string]interface{} {

	result := make( []map[string]interface{}, 0 )
	if model, err := models.GetModel("Product"); err == nil {
		if model, ok := model.(product.I_Product); ok {

			productsList, err := model.List()
			if err != nil { return jsonError(err) }

			for _, listValue := range productsList {
				if productItem, ok := listValue.(product.I_Product); ok {

					resultItem := map[string]interface{} {
						"_id": productItem.GetId(),
						"sku": productItem.GetSku(),
						"name": productItem.GetName(),
					}

					result = append(result, resultItem)
				}
			}

			return map[string]interface{} { "result": result }
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}

// WEB REST API used to create new one product
//   - product attributes must be included in POST form
//   - sku and name attributes required
func CreateProductRestAPI(req *http.Request, params map[string]string) map[string]interface{} {
	req.ParseForm()
	queryParams := req.PostForm

	if queryParams.Get("sku") == "" || queryParams.Get("name") == "" {
		return jsonError( errors.New("product 'name' and/or 'sku' was not specified") )
	}

	if model, err := models.GetModel("Product"); err == nil {
		if model, ok := model.(product.I_Product); ok {

			for attribute, value := range queryParams {
				err := model.Set(attribute, value[0])
				if err != nil { return jsonError(err) }
			}

			err := model.Save()
			if err != nil { return jsonError(err) }

			return model.ToHashMap()
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}


// WEB REST API used to update existing product
//   - product id must be specified in request URI
//   - product attributes must be included in POST form
func UpdateProductRestAPI(req *http.Request, params map[string]string) map[string]interface{} {

	//check request params
	productId, isSpecifiedId := params["id"]
	if !isSpecifiedId { return jsonError(errors.New("product 'id' was not specified")) }

	req.ParseForm()
	queryParams := req.PostForm
	if _, present := queryParams["_id"]; present { return jsonError(errors.New("_id attribute can't be updated")) }
	if len(queryParams) == 0 { return jsonError(errors.New("update attributes were not set")) }

	// update operations
	model, err := models.GetModel("Product")
	if err != nil { return jsonError(err) }

	productModel, ok := model.(product.I_Product)
	if !ok { return jsonError(errors.New("product type is not I_Product campatible")) }

	err = productModel.Load( productId )
	if err != nil { return jsonError(err) }

	for attrName, attrVal := range queryParams {
		err = productModel.Set(attrName, attrVal[0])
		if err != nil { return jsonError(err) }
	}

	err = productModel.Save()
	if err != nil { return jsonError(err) }

	return productModel.ToHashMap()
}
