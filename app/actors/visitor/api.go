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
func setupAPI() error {

	// Dashboard API
	err := api.GetRestService().RegisterAPI("visitor", "POST", "create", restCreateVisitor)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "PUT", "update", restUpdateVisitor)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "PUT", "update/:id", restUpdateVisitor)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "DELETE", "delete/:id", restDeleteVisitor)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("visitor", "GET", "load/:id", restGetVisitor)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "list", restListVisitors)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "attribute/list", restListVisitorAttributes)
	if err != nil {
		return err
	}

	// Storefront API
	err = api.GetRestService().RegisterAPI("visitor", "POST", "register", restRegister)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "validate/:key", restValidate)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "info", restInfo)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "GET", "logout", restLogout)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login", restLogin)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login-facebook", restLoginFacebook)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor", "POST", "login-google", restLoginGoogle)
	if err != nil {
		return err
	}
	/*err = api.GetRestService().RegisterAPI("visitor", "GET", "forgot-password", restForgotPassword)
	if err != nil {
		return err
	}*/

	return nil
}

// WEB REST API used to create new visitor
//   - visitor attributes must be included in POST form
//   - email attribute required
func restCreateVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email") {
		return nil, errors.New("'email' was not specified")
	}

	// create operation
	//-----------------
	visitorModel, err := visitor.GetVisitorModel()
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

// WEB REST API used to update existing visitor
//   - visitor id must be specified in request URI
//   - visitor attributes must be included in POST form
func restUpdateVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	visitorId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {

		sessionValue := params.Session.Get(visitor.SESSION_KEY_VISITOR_ID)
		sessionVisitorId, ok := sessionValue.(string)
		if !ok {
			return nil, errors.New("you are not logined in")
		}
		visitorId = sessionVisitorId
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// update operation
	//-----------------
	visitorModel, err := visitor.LoadVisitorById(visitorId)
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
func restDeleteVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	visitorId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor id was not specified")
	}

	// delete operation
	//-----------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, err
	}

	err = visitorModel.Delete(visitorId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func restGetVisitor(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	visitorId, isSpecifiedId := params.RequestURLParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor id was not specified")
	}

	// get operation
	//--------------
	visitorModel, err := visitor.LoadVisitorById(visitorId)
	if err != nil {
		return nil, err
	}

	return visitorModel.ToHashMap(), nil
}

// WEB REST API function used to get visitors list
func restListVisitors(params *api.T_APIHandlerParams) (interface{}, error) {
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, err
	}

	return visitorModel.List()
}

// WEB REST API function used to obtain visitor attributes information
func restListVisitorAttributes(params *api.T_APIHandlerParams) (interface{}, error) {
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, err
	}

	attrInfo := visitorModel.GetAttributesInfo()
	return attrInfo, nil
}

// WEB REST API used to register new visitor (same as create but with email validation)
//   - visitor attributes must be included in POST form
//   - email attribute required
func restRegister(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email") {
		return nil, errors.New("email was not specified")
	}

	// register visitor operation
	//---------------------------
	visitorModel, err := visitor.GetVisitorModel()
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

	visitorModel.Invalidate()

	return visitorModel.ToHashMap(), nil
}

func restValidate(params *api.T_APIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	validationKey, isKeySpecified := params.RequestURLParams["key"]
	if !isKeySpecified {
		return nil, errors.New("validation key was not specified")
	}

	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, err
	}

	err = visitorModel.Validate(validationKey)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API function used to obtain visitor information
//   - visitor id must be specified in request URI
func restInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	sessionValue := params.Session.Get("visitor_id")
	visitorId, ok := sessionValue.(string)
	if !ok {
		return "you are not logined in", nil
	}

	// visitor info
	//--------------
	visitorModel, err := visitor.LoadVisitorById(visitorId)
	if err != nil {
		return nil, err
	}

	result := visitorModel.ToHashMap()
	result["facebook_id"] = visitorModel.GetFacebookId()
	result["google_id"] = visitorModel.GetGoogleId()

	return result, nil
}

// WEB REST API function used to make visitor logout
func restLogout(params *api.T_APIHandlerParams) (interface{}, error) {
	sessionValue := params.Session.Get("visitor_id")
	if sessionValue != nil {
		return nil, errors.New("you are not logined in")
	}

	params.Session.Set(visitor.SESSION_KEY_VISITOR_ID, nil)

	return "ok", nil
}

// WEB REST API function used to make visitor login
//   - email and password information needed
func restLogin(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if !utils.KeysInMapAndNotBlank(reqData, "email", "password") {
		return nil, errors.New("email and/or password were not specified")
	}

	// visitor info
	//--------------
	visitorModel, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, err
	}

	err = visitorModel.LoadByEmail(utils.InterfaceToString(reqData["email"]))
	if err != nil {
		return nil, err
	}

	ok := visitorModel.CheckPassword(utils.InterfaceToString(reqData["password"]))
	if !ok {
		return nil, errors.New("wrong password")
	}

	if visitorModel.IsValidated() {
		params.Session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())
	} else {
		return nil, errors.New("visitor is not validated")
	}

	return "ok", nil
}

// WEB REST API function used to make login/registration via Facebook
//   - access_token and user_id params needed
//   - user needed information will be taken from Facebook
func restLoginFacebook(params *api.T_APIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if !utils.KeysInMapAndNotBlank(reqData, "access_token") {
		return nil, errors.New("access_token was not specified")
	}

	if !utils.KeysInMapAndNotBlank(reqData, "user_id") {
		return nil, errors.New("user_id was not specified")
	}

	// facebook login operation
	//-------------------------

	// using access token to get user information
	url := "https://graph.facebook.com/" + reqData["user_id"].(string) + "?access_token=" + reqData["access_token"].(string)
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
	err = json.Unmarshal(responseData, &jsonMap)
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
	err = visitorModel.LoadByFacebookId(reqData["user_id"].(string))
	if err != nil && strings.Contains(err.Error(), "not found") {

		// there is no such facebook_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail(jsonMap["email"].(string))
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

	params.Session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())

	return "ok", nil
}

// WEB REST API function used to make login/registration via Google
//   - access_token param needed
//   - user needed information will be taken from Google
func restLoginGoogle(params *api.T_APIHandlerParams) (interface{}, error) {
	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, errors.New("unexpected request content")
	}

	if !utils.KeysInMapAndNotBlank(reqData, "access_token") {
		return nil, errors.New("access_token was not specified")
	}

	// google login operation
	//-------------------------

	// using access token to get user information
	url := "https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + reqData["access_token"].(string)
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
	err = json.Unmarshal(responseData, &jsonMap)
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
	err = visitorModel.LoadByGoogleId(jsonMap["email"].(string))
	if err != nil && strings.Contains(err.Error(), "not found") {

		// there is no such google_id in DB, trying to find by e-mail
		err = visitorModel.LoadByEmail(jsonMap["email"].(string))
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

	params.Session.Set(visitor.SESSION_KEY_VISITOR_ID, visitorModel.GetId())

	return "ok", nil
}
