package visitor

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/api"
)

// REST API registration function
func (it *DefaultVisitor) setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor", "POST", "create", it.CreateVisitorRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "PUT", "update/:id", it.UpdateVisitorRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "DELETE", "delete/:id", it.DeleteVisitorRestAPI)
	if err != nil {
		return err
	}


	err = api.GetRestService().RegisterAPI("visitor", "GET", "load/:id", it.GetVisitorRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "list", it.ListVisitorsRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "attribute/list", it.ListVisitorAttributesRestAPI)
	if err != nil {
		return err
	}

	return nil
}



// WEB REST API used to create new visitor
//   - visitor attributes must be included in POST form
//   - email attribute required
func (it *DefaultVisitor) CreateVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//---------------------
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if queryParams["email"] == "" {
		return nil, errors.New("'email' was not specified")
	}


	// create operation
	//-----------------
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	for attribute, value := range queryParams {
		err := visitorModel.Set(attribute, value)
		if err != nil {
			return nil, err
		}
	}

	err = visitorModel.Save()
	if err != nil {
		return nil, err
	}

	return visitorModel.ToHashMap(), nil
}



// WEB REST API used to update existing visitor
//   - visitor id must be specified in request URI
//   - visitor attributes must be included in POST form
func (it *DefaultVisitor) UpdateVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//---------------------
	visitorId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor 'id' was not specified")
	}

	reqData, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}


	// update operation
	//-----------------
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	err = visitorModel.Load(visitorId)
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		err := visitorModel.Set(attribute, value)
		if err != nil {
			return nil, err
		}
	}

	err = visitorModel.Save()
	if err != nil {
		return nil, err
	}

	return visitorModel.ToHashMap(), nil
}



// WEB REST API used to delete visitor
//   - visitor id must be specified in request URI
func (it *DefaultVisitor) DeleteVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//---------------------
	visitorId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor 'id' was not specified")
	}


	// delete operation
	//-----------------
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	err = visitorModel.Delete(visitorId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}



// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func (it *DefaultVisitor) GetVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//---------------------
	visitorId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor 'id' was not specified")
	}

	// get operation
	//--------------
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	err = visitorModel.Load(visitorId)
	if err != nil {
		return nil, err
	}

	return visitorModel.ToHashMap(), nil
}



// WEB REST API function used to get visitors list
func (it *DefaultVisitor) ListVisitorsRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	result := make([]map[string]interface{}, 0)

	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	visitorsList, err := visitorModel.List()
	if err != nil {
		return nil, err
	}

	for _, listValue := range visitorsList {
		if visitorItem, ok := listValue.(visitor.I_Visitor); ok {

			resultItem := map[string]interface{} {
				"_id":           		visitorItem.GetId(),
				"email":           		visitorItem.GetEmail(),
				"full_name":           	visitorItem.GetFullName(),
			}
			result = append(result, resultItem)
		}
	}

	return result, nil
}



// WEB REST API function used to obtain visitor attributes information
func (it *DefaultVisitor) ListVisitorAttributesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	vis, isObject := model.(models.I_Object)
	if !isObject {
		return nil, errors.New("visitor model is not I_Object compatible")
	}

	attrInfo := vis.GetAttributesInfo()
	return attrInfo, nil
}
