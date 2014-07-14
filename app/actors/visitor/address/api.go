package address

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/api"
)


// REST API registration function
func (it *DefaultVisitorAddress) setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor/address", "POST", "create", it.CreateVisitorAddressRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "PUT", "update/:id", it.UpdateVisitorAddressRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "DELETE", "delete/:id", it.DeleteVisitorAddressRestAPI)
	if err != nil {
		return err
	}


	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "attribute/list", it.ListVisitorAddressAttributesRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "list/:visitorId", it.ListVisitorAddressRestAPI)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "load/:id", it.GetVisitorAddressRestAPI)
	if err != nil {
		return err
	}

	return nil
}



// WEB REST API used to create new visitor address
//   - visitor address attributes must be included in POST form
//   - visitor id required
func (it *DefaultVisitorAddress) CreateVisitorAddressRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//---------------------
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if _, ok := queryParams["visitor_id"]; !ok {
		return nil, errors.New("visitor id must be specified")
	}

	// create visitor address operation
	//---------------------------------
	model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil, err
	}

	addressModel, ok := model.(visitor.I_VisitorAddress)
	if !ok {
		return nil, errors.New("visitor address model is not I_VisitorAddress compatible")
	}

	for attribute, value := range queryParams {
		err := addressModel.Set(attribute, value)
		if err != nil {
			return nil, err
		}
	}

	err = addressModel.Save()
	if err != nil {
		return nil, err
	}

	return addressModel.ToHashMap(), nil
}



// WEB REST API used to update existing visitor address
//   - visitor address id must be specified in request URI
//   - visitor address attributes must be included in POST form
func (it *DefaultVisitorAddress) UpdateVisitorAddressRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//---------------------
	addressId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor address 'id' was not specified")
	}
	reqData, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}


	// update operation
	//-----------------
	model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil, err
	}

	addressModel, ok := model.(visitor.I_VisitorAddress)
	if !ok {
		return nil, errors.New("visitor address model is not I_Visitor campatible")
	}

	err = addressModel.Load(addressId)
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		err := addressModel.Set(attribute, value)
		if err != nil {
			return nil, err
		}
	}

	err = addressModel.Save()
	if err != nil {
		return nil, err
	}

	return addressModel.ToHashMap(), nil
}



// WEB REST API used to delete visitor address
//   - visitor address attributes must be included in POST form
func (it *DefaultVisitorAddress) DeleteVisitorAddressRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	// check request params
	//--------------------
	addressId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor address 'id' was not specified")
	}

	model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil, err
	}

	addressModel, ok := model.(visitor.I_VisitorAddress)
	if !ok {
		return nil, errors.New("visitor address model is not I_VisitorAddress campatible")
	}

	// delete operation
	//-----------------
	err = addressModel.Delete(addressId)
	if err != nil {
		return nil, err
	}

	return "ok", nil
}



// WEB REST API function used to obtain visitor address attributes information
func (it *DefaultVisitorAddress) ListVisitorAddressAttributesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
	model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil, err
	}

	address, isObject := model.(models.I_Object)
	if !isObject {
		return nil, errors.New("address address model is not I_Object compatible")
	}

	attrInfo := address.GetAttributesInfo()
	return attrInfo, nil
}



// WEB REST API function used to obtain visitor addresses list
//   - visitor id must be specified in request URI
func (it *DefaultVisitorAddress) ListVisitorAddressRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	result := make([]map[string]interface{}, 0)

	// check request params
	//---------------------
	visitorId, isSpecifiedId := reqParams["visitorId"]
	if !isSpecifiedId {
		return nil, errors.New("visitor 'id' was not specified")
	}

	// list operation
	//---------------
	model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil, err
	}

	addressModel, ok := model.(visitor.I_VisitorAddress)
	if !ok {
		return nil, errors.New("VisitorAddress model is not I_VisitorAddress compatible")
	}

	addressModel.ListFilterAdd("visitor_id", "=", visitorId)
	addressesList, err := addressModel.List()
	if err != nil {
		return nil, err
	}

	for _, listValue := range addressesList {
		if addressItem, ok := listValue.(visitor.I_VisitorAddress); ok {
			result = append(result, addressItem.ToHashMap())
		}
	}

	return result, nil
}



// WEB REST API used to get visitor address object
//   - visitor address id must be specified in request URI
func (it *DefaultVisitorAddress) GetVisitorAddressRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
	visitorId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor 'id' was not specified")
	}

	if model, err := models.GetModel("VisitorAddress"); err == nil {
		if model, ok := model.(visitor.I_VisitorAddress); ok {

			err = model.Load(visitorId)
			if err != nil {
				return nil, err
			}

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}
