package api

import (
	"net/http"
)

// interface for accessing private storage assigned to particular API request
type I_Session interface {
	GetId() string

	Get(key string) interface{}
	Set(key string, value interface{})

	Close() error
}

// interface to interact with RESTFul API service
type I_RestService interface {
	GetName() string

	Run() error
	RegisterAPI(service string, method string, uri string, handler F_APIHandler) error

	http.Handler
}
