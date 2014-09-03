package category

import (
	"errors"

	"strconv"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models/category"
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
	categoryModel, err := category.GetCategoryModel()
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

				categoryModel.ListLimit(offset, limit)
			} else if len(splitResult) > 0 {
				limit, err := strconv.Atoi(strings.TrimSpace(splitResult[0]))
				if err != nil {
					return nil, err
				}

				categoryModel.ListLimit(0, limit)
			} else {
				categoryModel.ListLimit(0, 0)
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
			err := categoryModel.ListAddExtraAttribute(strings.TrimSpace(extraAttribute))
			if err != nil {
				return nil, err
			}
		}
	}

	return categoryModel.List()
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

// WEB REST API function used to list product in category
//   - category id must be specified in request URI
func restListCategoryProducts(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

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

	// synthetic product list limit
	offset := 0
	limit := -1

	// limit parameter handler
	if limitParam, isLimit := reqData["limit"]; isLimit {
		if limitParam, ok := limitParam.(string); ok {
			splitResult := strings.Split(limitParam, ",")
			if len(splitResult) > 1 {

				offset, err = strconv.Atoi(strings.TrimSpace(splitResult[0]))
				if err != nil {
					return nil, err
				}

				limit, err = strconv.Atoi(strings.TrimSpace(splitResult[1]))
				if err != nil {
					return nil, err
				}
			} else if len(splitResult) > 0 {
				limit, err = strconv.Atoi(strings.TrimSpace(splitResult[0]))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	products := categoryModel.GetProducts()

	i := 0
	result := make([]map[string]interface{}, 0)
	for _, product := range products {
		if limit == 0 {
			break
		}
		if i >= offset {
			limit -= 1

			// preparing product information
			productInfo := product.ToHashMap()
			if defaultImage, present := productInfo["default_image"]; present {
				mediaPath, err := product.GetMediaPath("image")
				if defaultImage, ok := defaultImage.(string); ok && defaultImage != "" && err == nil {
					productInfo["default_image"] = mediaPath + defaultImage
				}
			}
			result = append(result, productInfo)
		}
		i += 1
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

	return len(categoryModel.GetProducts()), nil
}

// WEB REST API function used to categories menu
func restGetCategoriesTree(params *api.T_APIHandlerParams) (interface{}, error) {

	var result = make([]map[string]interface{}, 0)

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return nil, errors.New("can't get DB engine")
	}

	collection, err := dbEngine.GetCollection(CATEGORY_COLLECTION_NAME)
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
