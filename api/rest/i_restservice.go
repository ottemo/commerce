package rest

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"encoding/json"
	"encoding/xml"

	"github.com/julienschmidt/httprouter"
	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/api/session"
)

// returns implementation name of our REST API service
func (it *DefaultRestService) GetName() string {
	return "httprouter"
}

// modules should call this function in order to provide own REST API functionality
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

			// TODO: should be separated on 2 cases
			req.ParseMultipartForm(32 << 20) // 32 MB
			if req.MultipartForm != nil {
				for attribute, value := range req.MultipartForm.Value {
					newContent[attribute] = value[0]
				}
			}

			req.ParseForm()
			for attribute, value := range req.PostForm {
				newContent[attribute] = value[0]
			}

			content = newContent

		default:
			var body []byte = nil

			if req.ContentLength > 0 {
				body = make([]byte, req.ContentLength)
				req.Body.Read(body)
			} else {
				body, _ = ioutil.ReadAll(req.Body)
			}

			content = string(body)
		}

		// starting session for request
		session, err := session.StartSession(req, resp)
		if err != nil {
			log.Println("Session init fail: " + err.Error())
		}

		// module handler callback
		apiParams := new(api.T_APIHandlerParams)
		apiParams.Request = req
		apiParams.RequestURLParams = mappedParams
		apiParams.RequestContent = content
		apiParams.ResponseWriter = resp
		apiParams.Session = session

		result, err := handler(apiParams)
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
func (it DefaultRestService) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	// CORS fix-up
	responseWriter.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie, Content-Length, Accept-Encoding, X-CSRF-Token")

	if request.Method == "GET" || request.Method == "POST" || request.Method == "PUT" || request.Method == "DELETE" {

		// default output format
		responseWriter.Header().Set("Content-Type", "application/json")

		it.Router.ServeHTTP(responseWriter, request)
	}
}

// REST server startup function - makes it to "ListenAndServe"
func (it *DefaultRestService) Run() error {
	log.Println("REST API Service [HTTPRouter] starting to listen on " + it.ListenOn)
	log.Fatal(http.ListenAndServe(it.ListenOn, it))

	return nil
}
