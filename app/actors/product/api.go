package product

import (
	"mime"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

func setupAPI() error {

	var err error = nil

	// 1. DefaultProduct API
	//----------------------
	err = api.GetRestService().RegisterAPI("product", "GET", "get/:id", restGetProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "create", restCreateProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "PUT", "update/:id", restUpdateProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "delete/:id", restDeleteProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("product", "GET", "attribute/list", restListProductAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "attribute/remove/:attribute", restRemoveProductAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "attribute/add", restAddProductAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("product", "GET", "media/get/:productId/:mediaType/:mediaName", restMediaGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/list/:productId/:mediaType", restMediaList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/path/:productId/:mediaType", restMediaPath)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "media/add/:productId/:mediaType/:mediaName", restMediaAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "media/remove/:productId/:mediaType/:mediaName", restMediaRemove)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// 2. DefaultProductCollection API
	//--------------------------------
	err = api.GetRestService().RegisterAPI("product", "GET", "list", restListProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "list", restListProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "count", restCountProducts)
	if err != nil {
		return env.ErrorDispatch(err)
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
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := reqData["Attribute"]
	if !isSpecified {
		return nil, env.ErrorNew("attribute name was not specified")
	}

	attributeLabel, isSpecified := reqData["Label"]
	if !isSpecified {
		return nil, env.ErrorNew("attribute label was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// make product attribute operation
	//---------------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attribute := models.T_AttributeInfo{
		Model:      product.MODEL_NAME_PRODUCT,
		Collection: COLLECTION_NAME_PRODUCT,
		Attribute:  utils.InterfaceToString(attributeName),
		Type:       "text",
		IsRequired: false,
		IsStatic:   false,
		Label:      utils.InterfaceToString(attributeLabel),
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
		return nil, env.ErrorDispatch(err)
	}

	return attribute, nil
}

// WEB REST API function used to remove custom attribute of product
func restRemoveProductAttribute(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	attributeName, isSpecified := params.RequestURLParams["attribute"]
	if !isSpecified {
		return nil, env.ErrorNew("attribute name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// remove attribute actions
	//-------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.RemoveAttribute(attributeName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	// load product operation
	//-----------------------
	productModel, err := product.LoadProductById(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(params.RequestURLParams, "sku", "name") {
		return nil, env.ErrorNew("product name and/or sku were not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// create product operation
	//-------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		err := productModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = productModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// delete operation
	//-----------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorNew("unexpected request content")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// update operations
	//------------------
	productModel, err := product.LoadProductById(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attrName, attrVal := range reqData {
		err = productModel.Set(attrName, attrVal)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = productModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew("media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaList, err := productModel.GetMediaPath(mediaType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew("media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaList, err := productModel.ListMedia(mediaType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew("media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, env.ErrorNew("media name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// income file processing
	//-----------------------
	file, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	fileSize, _ := file.Seek(0, 2)
	fileContents := make([]byte, fileSize)

	file.Seek(0, 0)
	file.Read(fileContents)

	// add media operation
	//--------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.AddMedia(mediaType, mediaName, fileContents)
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew("media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, env.ErrorNew("media name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.RemoveMedia(mediaType, mediaName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorNew("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew("media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, env.ErrorNew("media name was not specified")
	}

	params.ResponseWriter.Header().Set("Content-Type", mime.TypeByExtension(mediaName))

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
		return nil, env.ErrorDispatch(err)
	}

	// api function routine
	//---------------------
	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
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
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return productCollectionModel.List()
}

// WEB REST API function used to obtain visitors count in model collection
func restCountProducts(params *api.T_APIHandlerParams) (interface{}, error) {

	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := productCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}
