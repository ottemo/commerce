package rest_service

import (
	"net/http"
)

type I_RestService interface {
	GetName() string

	Run() error
	RegisterJsonAPI(service string, uri string, handler func(req *http.Request) map[string]interface{} ) error
}
