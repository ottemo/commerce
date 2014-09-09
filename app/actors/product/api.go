package product

import (
	"errors"
	"mime"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

func setupAPI() error {

	var err error = nil

	// 1. DefaultProduct API
	//----------------------
	err = api.GetRestService().RegisterAPI("product", "GET", "get/:id", restGetProduct)
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

	// 2. DefaultProductCollection API
	//--------------------------------
	err = api.GetRestService().RegisterAPI("product", "GET", "list", restListProducts)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "list", restListProducts)
	if err != nil {
		return err
	}

	return nil
}

//----------------------
// 1. DefaultProduct API
//----------------------

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
		Validators: "",
		IsLayered:  false,
	}

	for key, value := range reqData {
		switch strings.ToLower(key) {
		case "type":
			attribute.Type = utils.InterfaceToString(value)
		case "group":
			attribute.Group = utils.InterfaceToString(value)
		case "editors":
			attribute.Editors = utils.InterfaceToString(value)
		case "options":
			attribute.Options = utils.InterfaceToString(value)
		case "default":
			attribute.Default = utils.InterfaceToString(value)
		case "validators":
			attribute.Validators = utils.InterfaceToString(value)
		case "isrequired", "required":
			attribute.IsRequired = utils.InterfaceToBool(value)
		case "islayered", "layered":
			attribute.IsLayered = utils.InterfaceToBool(value)
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
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}

	err = productModel.RemoveAttribute(attributeName)
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

//--------------------------------
// 2. DefaultProductCollection API
//--------------------------------

// WEB REST API function used to obtain product list we have in database
//   - only [_id, sku, name] attributes returns by default
func restListProducts(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// api function routine
	//---------------------
	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, err
	}

	// limit parameter handle
	productCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, productCollectionModel.GetDBCollection())

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := productCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, err
			}
		}
	}

	return productCollectionModel.List()
}
