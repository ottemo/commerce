package api

import (
	"net/http"
)

type T_APIHandlerParams struct {
	ResponseWriter http.ResponseWriter
	Request *http.Request
	RequestURLParams map[string]string
	RequestContent interface{}
	Session I_Session
}

type F_APIHandler func(params *T_APIHandlerParams) (interface{}, error)

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
