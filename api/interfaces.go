package api

import (
	"net/http"
)

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
