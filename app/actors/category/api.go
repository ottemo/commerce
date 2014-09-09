package category

import (
	"errors"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("category", "GET", "list", restListCategories)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "list", restListCategories)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "create", restCreateCategory)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "PUT", "update/:id", restUpdateCategory)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "DELETE", "delete/:id", restDeleteCategory)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "get/:id", restGetCategory)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("category", "GET", "products/:id", restListCategoryProducts)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/count/:id", restCategoryProductsCount)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "products/:id", restListCategoryProducts)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "layers/:id", restListCategoryLayers)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "layers/:id", restListCategoryLayers)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/add/:categoryId/:productId", restAddCategoryProduct)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/remove/:categoryId/:productId", restRemoveCategoryProduct)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("category", "GET", "attribute/list", restListCategoryAttributes)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("category", "GET", "tree", restGetCategoriesTree)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function used to obtain category list we have in database
//   - parent categories and categorys will not be present in list
func restListCategories(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, errors.New("unexpected request content")
		} else {
			reqData = make(map[string]interface{})
		}
	}

	// operation start
	//----------------
	categoryCollectionModel, err := category.GetCategoryCollectionModel()
	if err != nil {
		return nil, err
	}

	// limit parameter handler
	categoryCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, categoryCollectionModel.GetDBCollection())

	// extra parameter handler
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := categoryCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, err
			}
		}
	}

	return categoryCollectionModel.List()
}

// WEB REST API used to create new category
//   - category attributes must be included in POST form
//   - name attribute required
func restCreateCategory(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if _, present := reqData["name"]; !present {
		return nil, errors.New("category name was not specified")
	}

	// create category operation
	//-------------------------
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		err := categoryModel.Set(attribute, value)
		if err != nil {
			return nil, err
		}
	}

	err = categoryModel.Save()
	if err != nil {
		return nil, err
	}

	return categoryModel.ToHashMap(), nil
}

// WEB REST API used to delete category
func restDeleteCategory(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	categoryId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category id was not specified")
	}

	// delete operation
	//-----------------
	categoryModel, err := category.GetCategoryModelAndSetId(categoryId)
	if err != nil {
		return nil, err
	}

	err = categoryModel.Delete()
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API used to update existing category
//   - category id must be specified in request URI
//   - category attributes must be included in POST form
func restUpdateCategory(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// update operations
	//------------------
	categoryModel, err := category.LoadCategoryById(categoryId)
	if err != nil {
		return nil, err
	}

	for attrName, attrVal := range reqData {
		err = categoryModel.Set(attrName, attrVal)
		if err != nil {
			return nil, err
		}
	}

	err = categoryModel.Save()
	if err != nil {
		return nil, err
	}

	return categoryModel.ToHashMap(), nil
}

// WEB REST API function used to obtain category attributes information
func restListCategoryAttributes(params *api.T_APIHandlerParams) (interface{}, error) {
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, err
	}

	attrInfo := categoryModel.GetAttributesInfo()

	return attrInfo, nil
}

// WEB REST API function used to obtain layered navigation options for products in category
//   - category id must be specified in request URI
func restListCategoryLayers(params *api.T_APIHandlerParams) (interface{}, error) {
	categoryId := params.RequestURLParams["id"]

	categoryModel, err := category.LoadCategoryById(categoryId)
	if err != nil {
		return nil, err
	}

	productsCollection := categoryModel.GetProductsCollection()
	productsDBCollection := productsCollection.GetDBCollection()

	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, err
	}
	productAttributesInfo := productModel.GetAttributesInfo()

	result := make(map[string]interface{})

	api.ApplyFilters(params, productsDBCollection)

	for _, productAttribute := range productAttributesInfo {
		if productAttribute.Layered {
			distinctValues, _ := productsDBCollection.Distinct(productAttribute.Attribute)
			result[productAttribute.Attribute] = distinctValues
		}
	}

	return result, nil
}

// WEB REST API function used to list product in category
//   - category id must be specified in request URI
func restListCategoryProducts(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	/*reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}*/

	categoryId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category id was not specified")
	}

	// product list operation
	//-----------------------
	categoryModel, err := category.LoadCategoryById(categoryId)
	if err != nil {
		return nil, err
	}

	productsCollection := categoryModel.GetProductsCollection()

	api.ApplyFilters(params, productsCollection.GetDBCollection())

	// preparing product information
	result := make([]map[string]interface{}, 0)

	for _, product := range productsCollection.ListProducts() {
		productInfo := product.ToHashMap()
		if defaultImage, present := productInfo["default_image"]; present {
			mediaPath, err := product.GetMediaPath("image")
			if defaultImage, ok := defaultImage.(string); ok && defaultImage != "" && err == nil {
				productInfo["default_image"] = mediaPath + defaultImage
			}
		}
		result = append(result, productInfo)
	}

	return result, nil
}

// WEB REST API function used to add product in category
//   - category and product ids must be specified in request URI
func restAddCategoryProduct(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryId, isSpecifiedId := params.RequestURLParams["categoryId"]
	if !isSpecifiedId {
		return nil, errors.New("category id was not specified")
	}
	productId, isSpecifiedId := params.RequestURLParams["productId"]
	if !isSpecifiedId {
		return nil, errors.New("product id was not specified")
	}

	// category product add operation
	//-------------------------------
	categoryModel, err := category.GetCategoryModelAndSetId(categoryId)
	if err != nil {
		return nil, err
	}

	err = categoryModel.AddProduct(productId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API function used to remove product from category
//   - category and product ids must be specified in request URI
func restRemoveCategoryProduct(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryId, isSpecifiedId := params.RequestURLParams["categoryId"]
	if !isSpecifiedId {
		return nil, errors.New("category id was not specified")
	}
	productId, isSpecifiedId := params.RequestURLParams["productId"]
	if !isSpecifiedId {
		return nil, errors.New("product id was not specified")
	}

	// category product add operation
	//-------------------------------
	categoryModel, err := category.GetCategoryModelAndSetId(categoryId)
	if err != nil {
		return nil, err
	}

	err = categoryModel.RemoveProduct(productId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API function used to obtain all product attributes
//   - product id must be specified in request URI "http://[site:port]/product/get/:id"
func restGetCategory(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category id was not specified")
	}

	// load product operation
	//-----------------------
	categoryModel, err := category.LoadCategoryById(categoryId)
	if err != nil {
		return nil, err
	}

	return categoryModel.ToHashMap(), nil
}

// WEB REST API function used to list product in category
//   - category id must be specified in request URI
func restCategoryProductsCount(params *api.T_APIHandlerParams) (interface{}, error) {

	categoryId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category id was not specified")
	}

	// product list operation
	//-----------------------
	categoryModel, err := category.LoadCategoryById(categoryId)
	if err != nil {
		return nil, err
	}

	// count when we have filters (more complex and slow)
	if len(params.RequestGETParams) > 0 {
		productsDBCollection := categoryModel.GetProductsCollection().GetDBCollection()
		api.ApplyFilters(params, productsDBCollection)
		return productsDBCollection.Count()
	}

	return len(categoryModel.GetProductIds()), nil
}

// WEB REST API function used to categories menu
func restGetCategoriesTree(params *api.T_APIHandlerParams) (interface{}, error) {

	var result = make([]map[string]interface{}, 0)

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return nil, errors.New("can't get DB engine")
	}

	collection, err := dbEngine.GetCollection(COLLECTION_NAME_CATEGORY)
	if err != nil {
		return nil, err
	}

	err = collection.AddSort("path", false)
	if err != nil {
		return nil, err
	}

	rowData, err := collection.Load()
	if err != nil {
		return nil, err
	}

	categoryStack := make([]map[string]interface{}, 0)
	pathStack := make([]string, 0)

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
					return nil, errors.New("category tree builder internal error")
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
