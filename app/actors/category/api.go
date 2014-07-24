package category

import(
	"errors"
	"net/http"

	"strings"
	"strconv"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"

	"github.com/ottemo/foundation/db"
)

func (it *DefaultCategory) setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("category", "GET", "list", it.ListCategoriesRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "list", it.ListCategoriesRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "create", it.CreateCategoryRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "PUT", "update/:id", it.UpdateCategoryRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "DELETE", "delete/:id", it.DeleteCategoryRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "get/:id", it.GetCategoryRestAPI)
	if err != nil {
		return err
	}


	err = api.GetRestService().RegisterAPI("category", "GET", "products/:id", it.ListCategoryProductsRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/count/:id", it.ListCategoryProductsCountRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "products/:id", it.ListCategoryProductsRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/add/:categoryId/:productId", it.AddCategoryProductRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/remove/:categoryId/:productId", it.RemoveCategoryProductRestAPI)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("category", "GET", "attribute/list", it.ListCategoryAttributesRestAPI)
	if err != nil {
		return err
	}


	err = api.GetRestService().RegisterAPI("category", "GET", "tree", it.GetCategoriesTreeRestAPI)
	if err != nil {
		return err
	}


	return nil
}


// WEB REST API function used to obtain category list we have in database
//   - parent categories and categorys will not be present in list
func (it *DefaultCategory) ListCategoriesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := reqContent.(map[string]interface{})
	if !ok {
		if req.Method == "POST" {
			return nil, errors.New("unexpected request content")
		} else {
			reqData = make(map[string]interface{})
		}
	}

	// operation start
	//----------------
	model, err := models.GetModel("Category")
	if err != nil {
		return nil, errors.New("'Category' model not defined")
	}

	categoryModel, compatible := model.(category.I_Category)
	if !compatible  {
		return nil, errors.New("Category model is not I_Category compatible")
	}

	// limit parameter handler
	if limit, isLimit := reqData["limit"]; isLimit {
		if limit, ok := limit.(string); ok {
			splitResult := strings.Split(limit, ",")
			if len(splitResult) > 1 {

				offset, err := strconv.Atoi( strings.TrimSpace(splitResult[0]) )
				if err != nil {
					return nil, err
				}

				limit, err := strconv.Atoi( strings.TrimSpace(splitResult[1]) )
				if err != nil {
					return nil, err
				}

				categoryModel.ListLimit(offset, limit)
			} else if len(splitResult) > 0 {
				limit, err := strconv.Atoi( strings.TrimSpace(splitResult[0]) )
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
func (it *DefaultCategory) CreateCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if _, present := reqData["name"]; !present {
		return nil, errors.New("category 'name' was not specified")
	}

	// create category operation
	//-------------------------
	if model, err := models.GetModel("Category"); err == nil {
		if model, ok := model.(category.I_Category); ok {

			for attribute, value := range reqData {
				err := model.Set(attribute, value)
				if err != nil {
					return nil, err
				}
			}

			err := model.Save()
			if err != nil {
				return nil, err
			}

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}



// WEB REST API used to delete category
func (it *DefaultCategory) DeleteCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//--------------------
	categoryId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category 'id' was not specified")
	}

	model, err := models.GetModel("Category")
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(category.I_Category)
	if !ok {
		return nil, errors.New("category type is not I_Category campatible")
	}

	// delete operation
	//-----------------
	err = categoryModel.Delete(categoryId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}



// WEB REST API used to update existing category
//   - category id must be specified in request URI
//   - category attributes must be included in POST form
func (it *DefaultCategory) UpdateCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//---------------------
	categoryId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category 'id' was not specified")
	}

	reqData, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	// update operations
	//------------------
	model, err := models.GetModel("Category")
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(category.I_Category)
	if !ok {
		return nil, errors.New("category type is not I_Category campatible")
	}

	err = categoryModel.Load(categoryId)
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
func (it *DefaultCategory) ListCategoryAttributesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
	model, err := models.GetModel("Category")
	if err != nil {
		return nil, err
	}

	object, isObject := model.(models.I_Object)
	if !isObject {
		return nil, errors.New("category model is not I_Object compatible")
	}

	attrInfo := object.GetAttributesInfo()
	return attrInfo, nil
}



// WEB REST API function used to list product in category
//   - category id must be specified in request URI
func (it *DefaultCategory) ListCategoryProductsRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := reqContent.(map[string]interface{})
	if !ok {
		if req.Method == "POST" {
			return nil, errors.New("unexpected request content")
		} else {
			reqData = make(map[string]interface{})
		}
	}

	categoryId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category 'id' was not specified")
	}

	// product list operation
	//-----------------------
	model, err := models.GetModel("Category")
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(category.I_Category)
	if !ok {
		return nil, errors.New("category type is not I_Category campatible")
	}

	err = categoryModel.Load(categoryId)
	if err != nil {
		return nil, err
	}


	// synthetic product list limit
	offset := 0
	limit  := -1

	// limit parameter handler
	if limitParam, isLimit := reqData["limit"]; isLimit {
		if limitParam, ok := limitParam.(string); ok {
			splitResult := strings.Split(limitParam, ",")
			if len(splitResult) > 1 {

				offset, err = strconv.Atoi( strings.TrimSpace(splitResult[0]) )
				if err != nil {
					return nil, err
				}

				limit, err = strconv.Atoi( strings.TrimSpace(splitResult[1]) )
				if err != nil {
					return nil, err
				}
			} else if len(splitResult) > 0 {
				limit, err = strconv.Atoi( strings.TrimSpace(splitResult[0]) )
				if err != nil {
					return nil, err
				}
			}
		}
	}


	products := categoryModel.GetProducts()

	i := 0;
	result := make( []map[string]interface{}, 0)
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
func (it *DefaultCategory) AddCategoryProductRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//---------------------
	categoryId, isSpecifiedId := reqParams["categoryId"]
	if !isSpecifiedId {
		return nil, errors.New("category 'id' was not specified")
	}
	productId, isSpecifiedId := reqParams["productId"]
	if !isSpecifiedId {
		return nil, errors.New("product 'id' was not specified")
	}

	// category product add operation
	//-------------------------------
	model, err := models.GetModel("Category")
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(category.I_Category)
	if !ok {
		return nil, errors.New("category type is not I_Category campatible")
	}

	categoryModel.SetId(categoryId)
	err = categoryModel.AddProduct(productId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}



// WEB REST API function used to remove product from category
//   - category and product ids must be specified in request URI
func (it *DefaultCategory) RemoveCategoryProductRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {


	// check request params
	//---------------------
	categoryId, isSpecifiedId := reqParams["categoryId"]
	if !isSpecifiedId {
		return nil, errors.New("category 'id' was not specified")
	}
	productId, isSpecifiedId := reqParams["productId"]
	if !isSpecifiedId {
		return nil, errors.New("product 'id' was not specified")
	}

	// category product add operation
	//-------------------------------
	model, err := models.GetModel("Category")
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(category.I_Category)
	if !ok {
		return nil, errors.New("category type is not I_Category campatible")
	}

	categoryModel.SetId(categoryId)
	err = categoryModel.RemoveProduct(productId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API function used to obtain all product attributes
//   - product id must be specified in request URI "http://[site:port]/product/get/:id"
func (it *DefaultCategory) GetCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//---------------------
	categoryId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category 'id' was not specified")
	}

	// load product operation
	//-----------------------
	if model, err := models.GetModel("Category"); err == nil {
		if model, ok := model.(category.I_Category); ok {

			err = model.Load(categoryId)
			if err != nil {
				return nil, err
			}

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}


// WEB REST API function used to list product in category
//   - category id must be specified in request URI
func (it *DefaultCategory) ListCategoryProductsCountRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	categoryId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("category 'id' was not specified")
	}

	// product list operation
	//-----------------------
	model, err := models.GetModel("Category")
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(category.I_Category)
	if !ok {
		return nil, errors.New("category type is not I_Category campatible")
	}

	err = categoryModel.Load(categoryId)
	if err != nil {
		return nil, err
	}

	return len(categoryModel.GetProducts()), nil
}



// WEB REST API function used to categories menu
func (it *DefaultCategory) GetCategoriesTreeRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

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
		currentItem["id"]    = row["_id"]
		currentItem["name"]  = row["name"]
		currentItem["child"] = make( []map[string]interface{}, 0 )

		// calculating current path
		currentPath, ok := row["path"].(string)
		if !ok {
			continue
		}

		for idx:=len(pathStack)-1; idx >= 0; idx-- {
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
