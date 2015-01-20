package category

import (
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("categories", api.ConstRESTOperationGet, APICategoriesList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("categories/count", api.ConstRESTOperationGet, APICategoriesCount)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("categories/tree", api.ConstRESTOperationGet, APICategoriesTree)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("category", api.ConstRESTOperationCreate, APICategoryCreate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category/:categoryID", api.ConstRESTOperationUpdate, APICategoryUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category/:categoryID", api.ConstRESTOperationDelete, APICategoryDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category/:categoryID", api.ConstRESTOperationGet, APICategoryGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category/:categoryID/layers", api.ConstRESTOperationGet, APICategoryLayers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("categories/attributes", api.ConstRESTOperationGet, APICategoryAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("category/:categoryID/products", api.ConstRESTOperationGet, APICategoryProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category/:categoryID/products/count", api.ConstRESTOperationGet, APICategoryProductsCount)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category/:categoryID/product/:productID", api.ConstRESTOperationCreate, APICategoryProductAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category/:categoryID/product/:productID", api.ConstRESTOperationDelete, APICategoryProductRemove)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APICategoriesList returns list of existing categories
func APICategoriesList(context api.InterfaceApplicationContext) (interface{}, error) {

	// operation start
	//----------------
	categoryCollectionModel, err := category.GetCategoryCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handler
	categoryCollectionModel.ListLimit(api.GetListLimit(context))

	// filters handle
	api.ApplyFilters(context, categoryCollectionModel.GetDBCollection())

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		categoryCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// extra parameter handler
	api.ApplyExtraAttributes(context, categoryCollectionModel)

	return categoryCollectionModel.List()
}

// APICategoriesCount returns count of existing categories
func APICategoriesCount(context api.InterfaceApplicationContext) (interface{}, error) {
	categoryCollectionModel, err := category.GetCategoryCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := categoryCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(context, dbCollection)

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		categoryCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	return dbCollection.Count()
}

// APICategoryCreate creates a new category
//   - category attributes must be provided via content
//   - name attribute required
func APICategoryCreate(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	reqData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if _, present := reqData["name"]; !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "521b50d6-0d98-491a-8e3a-37678fbbccfe", "category name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// create category operation
	//-------------------------
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		err := categoryModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = categoryModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return categoryModel.ToHashMap(), nil
}

// APICategoryDelete removes category
//   - category id must be specified as "id" argument
func APICategoryDelete(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//--------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f1507b00-337e-4903-8244-5cf959dde3a4", "category id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// delete operation
	//-----------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APICategoryUpdate modifies existing category
//   - category id must be specified as "id" argument
//   - category attributes must be specified in content
func APICategoryUpdate(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "389975e7-611c-4d6c-8b4d-bca450f5f7e7", "category id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// update operations
	//------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attrName, attrVal := range reqData {
		err = categoryModel.Set(attrName, attrVal)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = categoryModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return categoryModel.ToHashMap(), nil
}

// APICategoryAttributes enumerates category attributes
func APICategoryAttributes(context api.InterfaceApplicationContext) (interface{}, error) {
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := categoryModel.GetAttributesInfo()

	return attrInfo, nil
}

// APICategoryAttributes enumerates category attributes and their possible values which is used for layered navigation
//   - category id should be specified in "id" argument
func APICategoryLayers(context api.InterfaceApplicationContext) (interface{}, error) {
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "389975e7-611c-4d6c-8b4d-bca450f5f7e7", "category id was not specified")
	}

	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(context) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d46dadf8-373a-4247-a81e-fbbe39a7fe74", "category is not available")
	}

	productsCollection := categoryModel.GetProductsCollection()
	productsDBCollection := productsCollection.GetDBCollection()

	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	productAttributesInfo := productModel.GetAttributesInfo()

	result := make(map[string]interface{})

	api.ApplyFilters(context, productsDBCollection)

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		productsDBCollection.AddFilter("enabled", "=", true)
	}

	for _, productAttribute := range productAttributesInfo {
		if productAttribute.IsLayered {
			distinctValues, _ := productsDBCollection.Distinct(productAttribute.Attribute)
			result[productAttribute.Attribute] = distinctValues
		}
	}

	return result, nil
}

// APICategoryProducts returns category related products
//   - category id should be specified in "categoryID" argument
func APICategoryProducts(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "12400cff-34fe-4cf5-ac6e-41625f9e3d5a", "category id was not specified")
	}

	// product list operation
	//-----------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(context) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a6f080d-dfa4-4f8c-8a0c-ec31cbe1cd87", "category is not available")
	}

	productsCollection := categoryModel.GetProductsCollection()

	api.ApplyFilters(context, productsCollection.GetDBCollection())

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		productsCollection.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// preparing product information
	var result []map[string]interface{}

	for _, productModel := range productsCollection.ListProducts() {
		productInfo := productModel.ToHashMap()
		if defaultImage, present := productInfo["default_image"]; present {
			mediaPath, err := productModel.GetMediaPath("image")
			if defaultImage, ok := defaultImage.(string); ok && defaultImage != "" && err == nil {
				productInfo["default_image"] = mediaPath + defaultImage
			}
		}
		result = append(result, productInfo)
	}

	return result, nil
}

// APICategoryProductAdd adds product to category
//   - category id and product id  should be specified in "categoryID" and "productID" arguments
func APICategoryProductAdd(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2fc40bdc-3c8e-4a9c-910b-cca62cda1b77", "category id was not specified")
	}
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "190a3d60-7769-4908-b383-80bc143128da", "product id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// category product add operation
	//-------------------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.AddProduct(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APICategoryProductRemove removes product from category
//   - category id and product id  should be specified in "categoryID" and "productID" arguments
func APICategoryProductRemove(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dc083799-bccb-48c8-bd56-5dcd0a0f6031", "category id was not specified")
	}
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9ffd2626-4192-4726-849e-c7a2416fab3a", "product id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// category product add operation
	//-------------------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.RemoveProduct(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APICategoryGet return category information
//   - category id should be specified in "id" argument
func APICategoryGet(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3c336fd7-1a18-4aea-9eb0-460d746f8dfa", "category id was not specified")
	}

	// load product operation
	//-----------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(context) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "80615e04-f43d-42a4-9482-39a5e7f8ccb7", "category is not available")
	}

	return categoryModel.ToHashMap(), nil
}

// APICategoryProductsCount returns count of products within category
func APICategoryProductsCount(context api.InterfaceApplicationContext) (interface{}, error) {

	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e1003839-d91c-4445-96c8-30e481b8347e", "category id was not specified")
	}

	// product list operation
	//-----------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(context) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9c2efb6a-cd28-4dd2-ba7f-5c1fe7c0df30", "category is not available")
	}

	productsDBCollection := categoryModel.GetProductsCollection().GetDBCollection()
	api.ApplyFilters(context, productsDBCollection)

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		productsDBCollection.AddFilter("enabled", "=", true)
	}

	return productsDBCollection.Count()
}

// APICategoriesTree returns categories parent/child relation map
func APICategoriesTree(context api.InterfaceApplicationContext) (interface{}, error) {

	var result = make([]map[string]interface{}, 0)

	collection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = collection.AddFilter("enabled", "=", true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = collection.AddSort("path", false)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	rowData, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var categoryStack []map[string]interface{}
	var pathStack []string

	for _, row := range rowData {

		currentItem := make(map[string]interface{})
		currentItem["id"] = row["_id"]
		currentItem["name"] = row["name"]
		currentItem["child"] = make([]map[string]interface{}, 0)

		// calculating current path
		currentPath, ok := row["path"].(string)
		if !ok {
			continue
		}

		for idx := len(pathStack) - 1; idx >= 0; idx-- {
			parentPath := pathStack[idx]

			// if we found parent
			if strings.Contains(currentPath, parentPath) {
				parent := categoryStack[idx]

				parentChild, ok := parent["child"].([]map[string]interface{})
				if !ok {
					return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c94e54a4-53dc-4ed2-bac4-7e9d93958765", "category tree builder internal error")
				}

				parent["child"] = append(parentChild, currentItem)

				break
			} else {
				pathStack = pathStack[0:idx]
				categoryStack = categoryStack[0:idx]
			}
		}

		if len(categoryStack) == 0 {
			result = append(result, currentItem)
		}

		categoryStack = append(categoryStack, currentItem)
		pathStack = append(pathStack, currentPath)

	}

	return result, nil
}
