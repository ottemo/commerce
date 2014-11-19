package product

import (
	"math/rand"
	"mime"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

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

	err = api.GetRestService().RegisterAPI("product", "GET", "media/get/:productID/:mediaType/:mediaName", restMediaGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/list/:productID/:mediaType", restMediaList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/path/:productID/:mediaType", restMediaPath)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "media/add/:productID/:mediaType/:mediaName", restMediaAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "media/remove/:productID/:mediaType/:mediaName", restMediaRemove)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "related/:productID", restRelatedList)
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
func restListProductAttributes(params *api.StructAPIHandlerParams) (interface{}, error) {
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := productModel.GetAttributesInfo()

	return attrInfo, nil
}

// WEB REST API function used to add new one custom attribute
func restAddProductAttribute(params *api.StructAPIHandlerParams) (interface{}, error) {

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

	attribute := models.StructAttributeInfo{
		Model:      product.ConstModelNameProduct,
		Collection: ConstCollectionNameProduct,
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
func restRemoveProductAttribute(params *api.StructAPIHandlerParams) (interface{}, error) {

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
func restGetProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("product id was not specified")
	}

	// load product operation
	//-----------------------
	productModel, err := product.LoadProductById(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return productModel.ToHashMap(), nil
}

// WEB REST API used to create new one product
//   - product attributes must be included in POST form
//   - sku and name attributes required
func restCreateProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

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
func restDeleteProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	productID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("product id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// delete operation
	//-----------------
	productModel, err := product.GetProductModelAndSetId(productID)
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
func restUpdateProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
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
	productModel, err := product.LoadProductById(productID)
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
func restMediaPath(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
		return nil, env.ErrorNew("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew("media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productID)
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
func restMediaList(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
		return nil, env.ErrorNew("product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew("media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetId(productID)
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
func restMediaAdd(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
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
	productModel, err := product.GetProductModelAndSetId(productID)
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
func restMediaRemove(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
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
	productModel, err := product.GetProductModelAndSetId(productID)
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
func restMediaGet(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
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
	productModel, err := product.GetProductModelAndSetId(productID)
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
func restListProducts(params *api.StructAPIHandlerParams) (interface{}, error) {

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

func restRelatedList(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isSpecifiedID := params.RequestURLParams["productID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew("product id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	count := 5
	if utils.InterfaceToInt(reqData["count"]) > 0 {
		count = utils.InterfaceToInt(reqData["count"])
	}

	// load product operation
	//-----------------------
	productModel, err := product.LoadProductById(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// result := make([]models.StructListItem, 0)
	var result []models.StructListItem

	relatedPids := utils.InterfaceToArray(productModel.Get("related_pids"))

	if len(relatedPids) > count {
		// indexes := make([]int, 0)
		var indexes []int
		for len(indexes) < count {

			new := rand.Intn(len(relatedPids))

			inArray := false
			for _, b := range indexes {
				if utils.InterfaceToInt(b) == new {
					inArray = true
				}
			}
			if !inArray {
				indexes = append(indexes, new)
			}
		}
		for _, index := range indexes {
			if productID := utils.InterfaceToString(relatedPids[index]); productID != "" {
				if productModel, err := product.LoadProductById(productID); err == nil {
					if err == nil {
						resultItem := new(models.StructListItem)

						mediaPath, err := productModel.GetMediaPath("image")
						if err != nil {
							return result, env.ErrorDispatch(err)
						}

						resultItem.Id = productModel.GetId()
						resultItem.Name = "[" + productModel.GetSku() + "] " + productModel.GetName()
						resultItem.Image = ""
						resultItem.Desc = productModel.GetShortDescription()

						if productModel.GetDefaultImage() != "" {
							resultItem.Image = mediaPath + productModel.GetDefaultImage()
						}

						if extra, isExtra := reqData["extra"]; isExtra {
							resultItem.Extra = make(map[string]interface{})
							extra := utils.Explode(utils.InterfaceToString(extra), ",")
							for _, value := range extra {
								resultItem.Extra[value] = productModel.Get(value)
							}
						}

						result = append(result, *resultItem)
					}
				}
			}
		}
	} else {
		for _, productID := range relatedPids {
			if productID == "" {
				continue
			}

			productModel, err := product.LoadProductById(utils.InterfaceToString(productID))
			if err == nil {
				resultItem := new(models.StructListItem)

				mediaPath, err := productModel.GetMediaPath("image")
				if err != nil {
					return result, env.ErrorDispatch(err)
				}

				resultItem.Id = productModel.GetId()
				resultItem.Name = "[" + productModel.GetSku() + "] " + productModel.GetName()
				resultItem.Image = ""
				resultItem.Desc = productModel.GetShortDescription()

				if productModel.GetDefaultImage() != "" {
					resultItem.Image = mediaPath + productModel.GetDefaultImage()
				}

				if extra, isExtra := reqData["extra"]; isExtra {
					resultItem.Extra = make(map[string]interface{})
					extra := utils.Explode(utils.InterfaceToString(extra), ",")
					for _, value := range extra {
						resultItem.Extra[value] = productModel.Get(value)
					}
				}

				result = append(result, *resultItem)
			}
		}
	}

	return result, nil
}

// WEB REST API function used to obtain visitors count in model collection
func restCountProducts(params *api.StructAPIHandlerParams) (interface{}, error) {

	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := productCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}
