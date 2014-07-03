package rest

import (
	"errors"
	"log"
	"net/http"

	"encoding/json"
	"encoding/xml"

	"github.com/julienschmidt/httprouter"
)

func (it *DefaultRestService) GetName() string {
	return "Negroni"
}

func (it *DefaultRestService) RegisterAPI(service string, method string, uri string, handler func(resp http.ResponseWriter, req *http.Request, params map[string]string) (interface{}, error) ) error {

	wrappedHandler := func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {

		mappedParams := make(map[string]string)
		for _, param := range params {
			mappedParams[param.Key] = param.Value
		}

		resp.Header().Set("Content-Type", "application/json")

		result, err := handler(resp, req, mappedParams)

		if result != nil {
			if _, ok := result.([]byte); !ok {
				if resp.Header().Get("Content-Type") == "application/json" {
					result, _ = json.Marshal(map[string]interface{} {"result": result, "error": err})
				}

				if resp.Header().Get("Content-Type") == "text/xml" {
					result, _ = xml.Marshal( result )
				}
			}

			resp.Write( result.([]byte) )
		}
	}

	path := "/" + service + "/" + uri

	switch method {
	case "GET":
		it.Router.GET(path, wrappedHandler)
	case "PUT":
		it.Router.PUT(path, wrappedHandler)
	case "POST":
		it.Router.POST(path, wrappedHandler)
	case "DELETE":
		it.Router.DELETE(path, wrappedHandler)
	default:
		return errors.New("unsupported method '" + method + "'")
	}

	key := path + " {" + method + "}"
	it.Handlers[key] = wrappedHandler

	return nil
}

func (it DefaultRestService) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")


	if req.Method == "GET" || req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		it.Router.ServeHTTP(resp, req)
	}
}

func (it *DefaultRestService) Run() error {
	log.Println("REST API Service [HTTPRouter] starting to listen on " + it.ListenOn)
	log.Fatal( http.ListenAndServe(it.ListenOn, it) )

	return nil
}
