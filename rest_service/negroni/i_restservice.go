package negroni

import (
	"errors"
	"net/http"
	"encoding/json"
)

func (it *NegroniRestService) GetName() string {
	return "Negroni"
}

func (it *NegroniRestService) RegisterJsonAPI(service string, uri string, handler func(req *http.Request) map[string]interface{} ) error {

	jsonHandler := func(resp http.ResponseWriter, req *http.Request) {
		result, _ := json.Marshal( handler(req) )

		resp.Header().Add("Content-Type", "application/json")
		resp.Write( result )
	}

	path := "/" + service + "/" + uri

	if _, present := it.Handlers[path]; present {
		return errors.New("There is already registered handler for " + path)
	}

	it.Mux.HandleFunc(path, jsonHandler)
	it.Handlers[path] = jsonHandler

	return nil
}

func (it *NegroniRestService) Run() error {
	it.Negroni.Run(it.ListenOn)

	return nil
}
