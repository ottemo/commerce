package category

import(
	"errors"
	"net/http"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
)

func (it *DefaultCategory) setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("category", "GET", "list", it.ListCategoriesRestAPI)
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


	return nil
}


// WEB REST API function used to obtain category list we have in database
//   - parent categories and categorys will not be present in list
func (it *DefaultCategory) ListCategoriesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	result := make([]map[string]interface{}, 0)

	model, err := models.GetModel("Category")
	if err != nil {
		return result, errors.New("'Category' model not defined")
	}

	categoryModel, compatible := model.(category.I_Category)
	if !compatible  {
		return result, errors.New("Category model is not I_Category compatible")
	}

	categoriesList, err := categoryModel.List()
	if err != nil {
		return nil, err
	}

	for _, listValue := range categoriesList {
		if categorytItem, ok := listValue.(category.I_Category); ok {
			result = append(result, categorytItem.ToHashMap() )
		}
	}

	return result, nil
}


// WEB REST API used to create new category
//   - category attributes must be included in POST form
//   - name attribute required
func (it *DefaultCategory) CreateCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

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
func (it *DefaultCategory) DeleteCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

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
func (it *DefaultCategory) UpdateCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

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
func (it *DefaultCategory) ListCategoryAttributesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
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
func (it *DefaultCategory) ListCategoryProductsRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//---------------------
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

	products := categoryModel.GetProducts()

	result := make( []map[string]interface{}, 0)
	for _, product := range products {
		result = append(result, product.ToHashMap())
	}


	return result, nil
}



// WEB REST API function used to add product in category
//   - category and product ids must be specified in request URI
func (it *DefaultCategory) AddCategoryProductRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

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
func (it *DefaultCategory) RemoveCategoryProductRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {


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
func (it *DefaultCategory) GetCategoryRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

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
