package default_visitor

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/visitor"
)

// sample: http://127.0.0.1:9000/visitor/load?id=5
func (it *DefaultVisitor) LoadVisitorAPI(resp http.ResponseWriter, req *http.Request, params map[string]string) (interface{}, error) {
	queryParams := req.URL.Query()

	productId := queryParams.Get("id")
	if productId == "" {
		return nil, errors.New("visitor 'id' was not specified")
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			err = model.Load( productId )
			if err != nil { return nil, err }

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}


// usage: http://127.0.0.1:9000/visitor/update?id=10&email=bad_guy@gmail.com&first_name=Bad&last_name=Guy
func (it *DefaultVisitor) UpdateVisitorAPI(resp http.ResponseWriter, req *http.Request, params map[string]string) (interface{}, error) {
	queryParams := req.URL.Query()

	if queryParams.Get("id") == "" {
		return nil, errors.New("visitor 'id' was not specified")
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			err := model.Load( queryParams.Get("id") )
			if err != nil { return nil, err }

			for attribute, value := range queryParams {
				err := model.Set(attribute, value[0])
				if err != nil { return nil, err }
			}

			err = model.Save()
			if err != nil { return nil, err }

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}

// usage: http://127.0.0.1:9000/visitor/create?email=bad_guy@gmail.com&first_name=Bad&last_name=Guy
func (it *DefaultVisitor) CreateVisitorAPI(resp http.ResponseWriter, req *http.Request, params map[string]string) (interface{}, error) {
	queryParams := req.URL.Query()

	if queryParams.Get("email") == "" {
		return nil, errors.New("'email' was not specified")
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			for attribute, value := range queryParams {
				err := model.Set(attribute, value[0])
				if err != nil { return nil, err }
			}

			err := model.Save()
			if err != nil { return nil, err }

			return model.ToHashMap(), nil
		}
	}

	return nil, errors.New("Something went wrong...")
}
