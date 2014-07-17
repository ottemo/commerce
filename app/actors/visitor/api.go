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

	// Dashboard API
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


	// Storefront API
	err = api.GetRestService().RegisterAPI("visitor", "POST", "register", it.RegisterRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "validate/:key", it.ValidateRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "info", it.InfoRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "logout", it.LogoutRestAPI)
	if err != nil {
		return err
	}
	/*err = api.GetRestService().RegisterAPI("visitor", "GET", "login", it.LoginRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "login-facebook", it.LoginRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "login-google", it.LoginRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "forget", it.ForgetPasswordRestAPI)
	if err != nil {
		return err
	}*/

	return nil
}



// WEB REST API used to create new visitor
//   - visitor attributes must be included in POST form
//   - email attribute required
func (it *DefaultVisitor) CreateVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

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
func (it *DefaultVisitor) UpdateVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

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
func (it *DefaultVisitor) DeleteVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

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
func (it *DefaultVisitor) GetVisitorRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

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
func (it *DefaultVisitor) ListVisitorsRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	return visitorModel.List()
}



// WEB REST API function used to obtain visitor attributes information
func (it *DefaultVisitor) ListVisitorAttributesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
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



// WEB REST API used to register new visitor (same as create but with email validation)
//   - visitor attributes must be included in POST form
//   - email attribute required
func (it *DefaultVisitor) RegisterRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	result, err := it.CreateVisitorRestAPI(resp, req, reqParams, reqContent, session)
	if err != nil {
		return result, err
	}

	// TODO: find better way to obtain customer id
	if result, ok := result.(map[string]interface{}); ok {
		if visitorId, ok := result["_id"]; ok {
			if visitorId != "" {
				session.Set("visitor_id", visitorId)
			}
		}
	}

	return result, nil
}


func (it *DefaultVisitor) ValidateRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
	// check request params
	//---------------------
	validationKey, isKeySpecified := reqParams["key"]
	if !isKeySpecified {
		return nil, errors.New("validation key was not specified")
	}

	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Object compatible")
	}

	visitorModel.Validate( validationKey )

	return validationKey, nil
}


// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func (it *DefaultVisitor) InfoRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	sessionValue := session.Get("visitor_id")
	visitorId, ok := sessionValue.(string)
	if !ok {
		return nil, errors.New("you are not logined in")
	}


	// visitor info
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


// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func (it *DefaultVisitor) LogoutRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
	sessionValue := session.Get("visitor_id")
	if sessionValue != nil {
		return nil, errors.New("you are not logined in")
	}

	session.Set("visitor_id", nil)

	return "ok", nil
}
