package api

import (
	"net/http"
)

type F_APIHandler func(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}, session I_Session) (interface{}, error)

type I_Session interface {
	GetId() string

	Get(key string) interface{}
	Set(key string, value interface{})
}

type I_RestService interface {
	GetName() string

	Run() error
	RegisterAPI(service string, method string, uri string, handler F_APIHandler) error

	http.Handler
}
