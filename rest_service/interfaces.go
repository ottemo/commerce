package rest_service

import (
	"net/http"
)

type I_RestService interface {
	GetName() string

	Run() error
	RegisterJsonAPI(service string, method string, uri string, handler func(req *http.Request, params map[string]string) map[string]interface{} ) error
}
