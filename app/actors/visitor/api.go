package visitor

import (
	"errors"
	"net/http"

	"io/ioutil"

	"encoding/json"
	"strings"

	"github.com/ottemo/foundation/app/utils"

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
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login", it.LoginRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login-facebook", it.LoginFacebookRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login-google", it.LoginGoogleRestAPI)
	if err != nil {
		return err
	}
	/*err = api.GetRestService().RegisterAPI("visitor", "GET", "forgot-password", it.ForgotPasswordRestAPI)
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

	// check request params
	//---------------------
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if queryParams["email"] == "" {
		return nil, errors.New("'email' was not specified")
	}


	// register visitor operation
	//---------------------------
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

	visitorModel.Invalidate()

	return visitorModel.ToHashMap(), nil
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

	err = visitorModel.Validate( validationKey )
	if err != nil {
		return nil, err
	}

	return "ok", nil
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

	result := visitorModel.ToHashMap()
	result["facebook_id"] = visitorModel.GetFacebookId()
	result["google_id"] = visitorModel.GetGoogleId()

	return result, nil
}


// WEB REST API function used to make visitor logout
func (it *DefaultVisitor) LogoutRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
	sessionValue := session.Get("visitor_id")
	if sessionValue != nil {
		return nil, errors.New("you are not logined in")
	}

	session.Set(visitor.SESSION_KEY_VISITOR_ID, nil)

	return "ok", nil
}


// WEB REST API function used to make visitor login
//   - email and password information needed
func (it *DefaultVisitor) LoginRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {

	// check request params
	//---------------------
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if _, ok := queryParams["email"].(string); !ok || queryParams["email"] == "" {
		return nil, errors.New("email was not specified")
	}

	if _, ok := queryParams["password"].(string); !ok || queryParams["password"] == "" {
		return nil, errors.New("password was not specified")
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

	err = visitorModel.LoadByEmail( queryParams["email"].(string) )
	if err != nil {
		return nil, err
	}

	ok = visitorModel.CheckPassword( queryParams["password"].(string) )
	if !ok {
		return nil, errors.New("wrong password")
	}

	if visitorModel.IsValidated() {
		session.Set(visitor.SESSION_KEY_VISITOR_ID , visitorModel.GetId())
	} else {
		return nil, errors.New("visitor is not validated")
	}

	return "ok", nil
}


// WEB REST API function used to make login/registration via Facebook
//   - access_token and user_id params needed
//   - user needed information will be taken from Facebook
func (it *DefaultVisitor) LoginFacebookRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
	// check request params
	//---------------------
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if _, ok := queryParams["access_token"].(string); !ok || queryParams["access_token"] == "" {
		return nil, errors.New("access_token was not specified")
	}

	if _, ok := queryParams["user_id"].(string); !ok ||  queryParams["user_id"] == "" {
		return nil, errors.New("user_id was not specified")
	}

	// facebook login operation
	//-------------------------

	// using access token to get user information
	url := "https://graph.facebook.com/" + queryParams["user_id"].(string) + "?access_token=" + queryParams["access_token"].(string)
	facebookResponse, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if facebookResponse.StatusCode != 200 {
		return nil, errors.New("Can't use google API: " + facebookResponse.Status)
	}

	var responseData []byte
	if facebookResponse.ContentLength > 0 {
		responseData = make([]byte, facebookResponse.ContentLength)
		_, err := facebookResponse.Body.Read(responseData)
		if err != nil {
			return nil, err
		}
	} else {
		responseData, err = ioutil.ReadAll(facebookResponse.Body)
		if err != nil {
			return nil, err
		}
	}

	// response json workaround
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(responseData , &jsonMap)
	if err != nil {
		return nil, err
	}

	if !utils.StrKeysInMap(jsonMap, "id", "email", "first_name", "last_name", "verified") {
		return nil, errors.New("unexpected facebook response")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	// trying to load visitor by facebook_id
	err = visitorModel.LoadByFacebookId( queryParams["user_id"].(string) )
	if err != nil && strings.Contains(err.Error(), "not found") {

		// there is no such facebook_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail( jsonMap["email"].(string) )
		if err != nil && strings.Contains(err.Error(), "not found") {
			// visitor not exists in out DB - reating new one
			visitorModel.Set("email", jsonMap["email"])
			visitorModel.Set("first_name", jsonMap["first_name"])
			visitorModel.Set("last_name", jsonMap["last_name"])
			visitorModel.Set("facebook_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, err
			}
		} else {
			// we have visitor with that e-mail just updating token, if it verified
			if value, ok := jsonMap["verified"].(bool); !(ok && value) {
				return nil, errors.New("facebook account email unverified")
			}

			visitorModel.Set("facebook_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, err
			}
		}
	}

	session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())

	return "ok", nil
}


// WEB REST API function used to make login/registration via Google
//   - access_token param needed
//   - user needed information will be taken from Google
func (it *DefaultVisitor) LoginGoogleRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session api.I_Session) (interface{}, error) {
	// check request params
	//---------------------
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if _, ok := queryParams["access_token"].(string); !ok || queryParams["access_token"] == "" {
		return nil, errors.New("access_token was not specified")
	}

	// google login operation
	//-------------------------

	// using access token to get user information
	url := "https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + queryParams["access_token"].(string)
	googleResponse, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if googleResponse.StatusCode != 200 {
		return nil, errors.New("Can't use google API: " + googleResponse.Status)
	}

	var responseData []byte
	if googleResponse.ContentLength > 0 {
		responseData = make([]byte, googleResponse.ContentLength)
		_, err := googleResponse.Body.Read(responseData)
		if err != nil {
			return nil, err
		}
	} else {
		responseData, err = ioutil.ReadAll(googleResponse.Body)
		if err != nil {
			return nil, err
		}
	}

	// response json workaround
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(responseData , &jsonMap)
	if err != nil {
		return nil, err
	}

	if !utils.StrKeysInMap(jsonMap, "id", "email", "verified_email", "given_name", "family_name") {
		return nil, errors.New("unexpected google response")
	}

	// trying to load visitor from our DB
	model, err := models.GetModel("Visitor")
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(visitor.I_Visitor)
	if !ok {
		return nil, errors.New("visitor model is not I_Visitor campatible")
	}

	// trying to load visitor by google_id
	err = visitorModel.LoadByGoogleId( jsonMap["email"].(string) )
	if err != nil && strings.Contains(err.Error(), "not found") {

		// there is no such google_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail( jsonMap["email"].(string) )
		if err != nil && strings.Contains(err.Error(), "not found") {

			// visitor e-mail not exists in out DB - creating new one
			visitorModel.Set("email", jsonMap["email"])
			visitorModel.Set("first_name", jsonMap["given_name"])
			visitorModel.Set("last_name", jsonMap["family_name"])
			visitorModel.Set("google_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, err
			}

		} else {
			// we have visitor with that e-mail just updating token, if it verified
			if value, ok := jsonMap["verified_email"].(bool); !(ok && value) {
				return nil, errors.New("google account email unverified")
			}

			visitorModel.Set("google_id", jsonMap["id"])

			err := visitorModel.Save()
			if err != nil {
				return nil, err
			}
		}
	}

	session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())

	return "ok", nil
}
