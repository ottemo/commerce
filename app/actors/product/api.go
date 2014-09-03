package product

import (
	"errors"
	"mime"
	"strconv"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("product", "GET", "list", restListProducts)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "get/:id", restGetProduct)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "list", restListProducts)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "create", restCreateProduct)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "PUT", "update/:id", restUpdateProduct)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "delete/:id", restDeleteProduct)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("product", "GET", "attribute/list", restListProductAttributes)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "attribute/remove/:attribute", restRemoveProductAttribute)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "attribute/add", restAddProductAttribute)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("product", "GET", "media/get/:productId/:mediaType/:mediaName", restMediaGet)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/list/:productId/:mediaType", restMediaList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/path/:productId/:mediaType", restMediaPath)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "media/add/:productId/:mediaType/:mediaName", restMediaAdd)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "media/remove/:productId/:mediaType/:mediaName", restMediaRemove)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function used to obtain product attributes information
func restListProductAttributes(params *api.T_APIHandlerParams) (interface{}, error) {
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}

	attrInfo := productModel.GetAttributesInfo()

	return attrInfo, nil
}

// WEB REST API function used to add new one custom attribute
func restAddProductAttribute(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	attributeName, isSpecified := reqData["Attribute"]
	if !isSpecified {
		return nil, errors.New("attribute name was not specified")
	}

	attributeLabel, isSpecified := reqData["Label"]
	if !isSpecified {
		return nil, errors.New("attribute label was not specified")
	}

	// make product attribute operation
	//---------------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}

	attribute := models.T_AttributeInfo{
		Model:      "product",
		Collection: "product",
		Attribute:  attributeName.(string),
		Type:       "text",
		IsRequired: false,
		IsStatic:   false,
		Label:      attributeLabel.(string),
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
	}

	for key, value := range reqData {
		switch value := value.(type) {
		case string:
			switch strings.ToLower(key) {
			case "type":
				attribute.Type = value
			case "required":
				if strings.ToLower(value) == "true" {
					attribute.IsRequired = true
				}
			case "group":
				attribute.Group = value
			case "editors":
				attribute.Editors = value
			case "options":
				attribute.Options = value
			case "default":
				attribute.Default = value
			}
		case bool:
			switch key {
			case "required":
				attribute.IsRequired = value
			}
		}
	}

	err = productModel.AddNewAttribute(attribute)
	if err != nil {
		return nil, err
	}

	return attribute, nil
}

// WEB REST API function used to remove custom attribute of product
func restRemoveProductAttribute(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	attributeName, isSpecified := params.RequestURLParams["attribute"]
	if !isSpecified {
		return nil, errors.New("attribute name was not specified")
	}

	// remove attribute actions
	//-------------------------
	produtcModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}

	err = produtcModel.RemoveAttribute(attributeName)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API function used to obtain all product attributes
//   - product id must be specified in request URI "http://[site:port]/product/get/:id"
func restGetProduct(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("product id was not specified")
	}

	// load product operation
	//-----------------------
	productModel, err := product.LoadProductById(productId)
	if err != nil {
		return nil, err
	}

	return productModel.ToHashMap(), nil
}

// WEB REST API function used to obtain product list we have in database
//   - only [_id, sku, name] attributes returns by default
func restEnumProducts(params *api.T_APIHandlerParams) (interface{}, error) {

	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}

	return productModel.List()
}

// WEB REST API function used to obtain product list we have in database
//   - only [_id, sku, name] attributes returns by default
func restListProducts(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// operation start
	//----------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}

	// limit parameter handler
	if limit, isLimit := reqData["limit"]; isLimit {
		if limit, ok := limit.(string); ok {
			splitResult := strings.Split(limit, ",")
			if len(splitResult) > 1 {

				offset, err := strconv.Atoi(strings.TrimSpace(splitResult[0]))
				if err != nil {
					return nil, err
				}

				limit, err := strconv.Atoi(strings.TrimSpace(splitResult[1]))
				if err != nil {
					return nil, err
				}

				err = productModel.ListLimit(offset, limit)
			} else if len(splitResult) > 0 {
				limit, err := strconv.Atoi(strings.TrimSpace(splitResult[0]))
				if err != nil {
					return nil, err
				}

				productModel.ListLimit(0, limit)
			} else {
				productModel.ListLimit(0, 0)
			}
		}
	}

	// extra parameter handler
	if extra, isExtra := reqData["extra"]; isExtra {
		extra, ok := extra.(string)
		if !ok {
			return nil, errors.New("extra parameter should be string")
		}

		splitResult := strings.Split(extra, ",")
		for _, extraAttribute := range splitResult {
			err := productModel.ListAddExtraAttribute(strings.TrimSpace(extraAttribute))
			if err != nil {
				return nil, err
			}
		}
	}

	return productModel.List()
}

// WEB REST API used to create new one product
//   - product attributes must be included in POST form
//   - sku and name attributes required
func restCreateProduct(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if !utils.KeysInMapAndNotBlank(params.RequestURLParams, "sku", "name") {
		return nil, errors.New("product name and/or sku were not specified")
	}

	// create product operation
	//-------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		err := productModel.Set(attribute, value)
		if err != nil {
			return nil, err
		}
	}

	err = productModel.Save()
	if err != nil {
		return nil, err
	}

	return productModel.ToHashMap(), nil
}

// WEB REST API used to delete product
//   - product attributes must be included in POST form
func restDeleteProduct(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	productId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("product id was not specified")
	}

	// delete operation
	//-----------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, err
	}

	err = productModel.Delete()
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API used to update existing product
//   - product id must be specified in request URI
//   - product attributes must be included in POST form
func restUpdateProduct(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("product id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, errors.New("unexpected request content")
	}

	// update operations
	//------------------
	productModel, err := product.LoadProductById(productId)
	if err != nil {
		return nil, err
	}

	for attrName, attrVal := range reqData {
		err = productModel.Set(attrName, attrVal)
		if err != nil {
			return nil, err
		}
	}

	err = productModel.Save()
	if err != nil {
		return nil, err
	}

	return productModel.ToHashMap(), nil
}

// WEB REST API used to add media for a product
//   - product id, media type must be specified in request URI
func restMediaPath(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productId, isIdSpecified := params.RequestURLParams["productId"]
	if !isIdSpecified {
		return nil, errors.New("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, errors.New("media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, err
	}

	mediaList, err := productModel.GetMediaPath(mediaType)
	if err != nil {
		return nil, err
	}

	return mediaList, nil
}

// WEB REST API used to add media for a product
//   - product id, media type must be specified in request URI
func restMediaList(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productId, isIdSpecified := params.RequestURLParams["productId"]
	if !isIdSpecified {
		return nil, errors.New("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, errors.New("media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, err
	}

	mediaList, err := productModel.ListMedia(mediaType)
	if err != nil {
		return nil, err
	}

	return mediaList, nil
}

// WEB REST API used to add media for a product
//   - product id, media type and media name must be specified in request URI
//   - media contents must be included as file in POST form
func restMediaAdd(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productId, isIdSpecified := params.RequestURLParams["productId"]
	if !isIdSpecified {
		return nil, errors.New("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, errors.New("media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, errors.New("media name was not specified")
	}

	// income file processing
	//-----------------------
	file, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, err
	}

	fileSize, _ := file.Seek(0, 2)
	fileContents := make([]byte, fileSize)

	file.Seek(0, 0)
	file.Read(fileContents)

	// add media operation
	//--------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, err
	}

	err = productModel.AddMedia(mediaType, mediaName, fileContents)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API used to add media for a product
//   - product id, media type and media name must be specified in request URI
func restMediaRemove(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productId, isIdSpecified := params.RequestURLParams["productId"]
	if !isIdSpecified {
		return nil, errors.New("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, errors.New("media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, errors.New("media name was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, err
	}

	err = productModel.RemoveMedia(mediaType, mediaName)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API used to get media contents for a product
//   - product id, media type and media name must be specified in request URI
func restMediaGet(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productId, isIdSpecified := params.RequestURLParams["productId"]
	if !isIdSpecified {
		return nil, errors.New("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, errors.New("media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, errors.New("media name was not specified")
	}

	params.ResponseWriter.Header().Set("Content-Type", mime.TypeByExtension(mediaName))

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, err
	}

	return productModel.GetMedia(mediaType, mediaName)
}
