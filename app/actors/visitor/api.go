package visitor

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/api"
)


func (it *DefaultVisitor) setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor", "POST", "create", it.CreateVisitorAPI)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("visitor", "PUT", "update/:id", it.UpdateVisitorAPI)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("visitor", "GET", "load/:id", it.LoadVisitorAPI)
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

// WEB REST API function used to obtain product attributes information
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

func (it *DefaultVisitor) ListVisitorsRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {

	result := make([]map[string]interface{}, 0)
	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			visitorsList, err := model.List()
			if err != nil {
				return nil, err
			}

			for _, listValue := range visitorsList {
				if visitorItem, ok := listValue.(visitor.I_Visitor); ok {

					resultItem := map[string]interface{}{
						"_id":           		visitorItem.GetId(),
						"email":           		visitorItem.GetEmail(),
						"full_name":           	visitorItem.GetFullName(),
						"first_name":          	visitorItem.GetFirstName(),
						"last_name":          	visitorItem.GetLastName(),
						"shipping_address":   	visitorItem.GetShippingAddress(),
						"billing_address": 		visitorItem.GetBillingAddress(),
					}


					result = append(result, resultItem)
				}
			}

			return result, nil
		}
	}

	return nil, errors.New("Something went wrong...")
}


// sample: http://127.0.0.1:9000/visitor/load?id=5
func (it *DefaultVisitor) LoadVisitorAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
	visitorId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor 'id' was not specified")
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			err = model.Load(visitorId)
			if err != nil {
				return nil, err
			}

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}

// usage: http://127.0.0.1:9000/visitor/update?id=10&email=bad_guy@gmail.com&first_name=Bad&last_name=Guy
func (it *DefaultVisitor) UpdateVisitorAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
	visitorId, isSpecifiedId := reqParams["id"]
	if !isSpecifiedId {
		return nil, errors.New("visitor 'id' was not specified")
	}
	reqData, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}


	if visitorId == "" {
		return nil, errors.New("visitor 'id' was not specified")
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			err := model.Load(visitorId)
			if err != nil {
				return nil, err
			}

			for attribute, value := range reqData {
				err := model.Set(attribute, value)
				if err != nil {
					return nil, err
				}
			}

			err = model.Save()
			if err != nil {
				return nil, err
			}

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}

// usage: http://127.0.0.1:9000/visitor/create?email=bad_guy@gmail.com&first_name=Bad&last_name=Guy
func (it *DefaultVisitor) CreateVisitorAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
//	queryParams := req.URL.Query()
	queryParams, ok := reqContent.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected request content")
	}

	if queryParams["email"] == "" {
		return nil, errors.New("'email' was not specified")
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

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
