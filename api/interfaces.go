package api

import (
	"net/http"
)

type I_RestService interface {
	GetName() string

	Run() error
	RegisterAPI(service string, method string, uri string, handler func(resp http.ResponseWriter, req *http.Request, params map[string]string) (interface{}, error) ) error

	http.Handler
}
