package address

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/api"
)

func (it *DefaultVisitorAddress) setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor/address", "GET", "attribute/list", it.ListAddressAttributesRestAPI)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("visitor/address", "POST", "create", it.CreateAddressAPI)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("visitor/address", "GET", "list/:visitorId", it.ListAddressRestAPI)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function used to obtain product attributes information
func (it *DefaultVisitorAddress) ListAddressAttributesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
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

func (it *DefaultVisitorAddress) CreateAddressAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
	//	queryParams := req.URL.Query()
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if model, err := models.GetModel("VisitorAddress"); err == nil {
		if model, ok := model.(visitor.I_VisitorAddress); ok {

			for attribute, value := range queryParams {
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

func (it *DefaultVisitorAddress) ListAddressRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	result := make([]map[string]interface{}, 0)
	if model, err := models.GetModel("VisitorAddress"); err == nil {
		if model, ok := model.(visitor.I_VisitorAddress); ok {

			addressesList, err := model.List()
			if err != nil {
				return nil, err
			}

			for _, listValue := range addressesList {
				if addressItem, ok := listValue.(visitor.I_VisitorAddress); ok {

					resultItem := map[string]interface{}{
						"_id":			addressItem.GetId(),
						"visitor_id":	addressItem.GetVisitorId(),
						"street":		addressItem.GetStreet(),
						"city":			addressItem.GetCity(),
						"state":		addressItem.GetState(),
						"phone":		addressItem.GetPhone(),
						"zip_code":		addressItem.GetZipCode(),
					}


					result = append(result, resultItem)
				}
			}

			return result, nil
		}
	}

	return nil, errors.New("Something went wrong...")
}

