package api

import (
	"net/http"
)

// interface for accessing private storage assigned to particular API request
type InterfaceSession interface {
	GetID() string

	Get(key string) interface{}
	Set(key string, value interface{})

	Close() error
}

// interface to interact with RESTFul API service
type InterfaceRestService interface {
	GetName() string

	Run() error
	RegisterAPI(service string, method string, uri string, handler FuncAPIHandler) error

	http.Handler
}
