package rest

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"encoding/json"
	"encoding/xml"

	"github.com/julienschmidt/httprouter"
	"github.com/ottemo/foundation/api"
)

// returns implementation name of our REST API service
func (it *DefaultRestService) GetName() string {
	return "httprouter"
}

// other modules should call this function in order to provide own REST API functionality
func (it *DefaultRestService) RegisterAPI(service string, method string, uri string, handler api.F_APIHandler) error {

	// httprouter needs other type of handler that we using
	wrappedHandler := func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {

		// getting URL params of request
		mappedParams := make(map[string]string)
		for _, param := range params {
			mappedParams[param.Key] = param.Value
		}

		// request content conversion (if possible)
		var content interface{} = nil

		contentType := req.Header.Get("Content-Type")
		switch {
		// JSON content
		case strings.Contains(contentType, "json"):
			newContent := map[string]interface{}{}

			buf := make([]byte, req.ContentLength)
			req.Body.Read(buf)
			json.Unmarshal(buf, &newContent)

			content = newContent

		// POST form content
		case strings.Contains(contentType, "form-data"):
			newContent := map[string]interface{}{}

			req.ParseForm()
			for attribute, value := range req.PostForm {
				newContent[attribute] = value
			}

			content = newContent
		}

		// module handler callback
		result, err := handler(resp, req, mappedParams, content)
		if err != nil {
			log.Printf("REST error: %s - %s\n", req.RequestURI, err.Error())
		}

		// result conversion before output
		if result != nil || err != nil {
			if _, ok := result.([]byte); !ok {

				// JSON encode
				if resp.Header().Get("Content-Type") == "application/json" {
					errorMsg := ""
					if err != nil {
						errorMsg = err.Error()
					}

					result, _ = json.Marshal(map[string]interface{}{"result": result, "error": errorMsg})
				}

				// XML encode
				if resp.Header().Get("Content-Type") == "text/xml" {
					result, _ = xml.Marshal(result)
				}
			}

			resp.Write(result.([]byte))
		}
	}

	// registration to httprouter
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

// entry point for HTTP request - takes control before request handled
// (go lang "http.server" package "Handler" interface implementation)
func (it DefaultRestService) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	// CORS fix-up
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")

	if req.Method == "GET" || req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {

		// default output format
		resp.Header().Set("Content-Type", "application/json")

		it.Router.ServeHTTP(resp, req)
	}
}

// REST server startup function - makes it to "ListenAndServe"
func (it *DefaultRestService) Run() error {
	log.Println("REST API Service [HTTPRouter] starting to listen on " + it.ListenOn)
	log.Fatal(http.ListenAndServe(it.ListenOn, it))

	return nil
}
