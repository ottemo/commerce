package json

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (it *DefaultService) GetName() string {
	return "Negroni"
}

func (it *DefaultService) RegisterJsonAPI(service string, method string, uri string, handler func(req *http.Request, params map[string]string) map[string]interface{}) error {

	jsonHandler := func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {

		mappedParams := make(map[string]string)
		for _, param := range params {
			mappedParams[param.Key] = param.Value
		}

		result, _ := json.Marshal(handler(req, mappedParams))

		resp.Header().Add("Content-Type", "application/json")
		resp.Write(result)
	}

	path := "/" + service + "/" + uri

	switch method {
	case "GET":
		it.Router.GET(path, jsonHandler)
	case "PUT":
		it.Router.PUT(path, jsonHandler)
	case "POST":
		it.Router.POST(path, jsonHandler)
	case "DELETE":
		it.Router.DELETE(path, jsonHandler)
	default:
		return errors.New("unsupported method '" + method + "'")
	}

	it.Handlers[path] = jsonHandler

	return nil
}

func (it *DefaultService) Run() error {
	log.Println("REST API Service [HTTPRouter] starting to listen on " + it.ListenOn)
	log.Fatal(http.ListenAndServe(it.ListenOn, it.Router))

	return nil
}
