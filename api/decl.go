package api

import (
	"net/http"
)

var (
	SESSION_KEY_ADMIN_RIGHTS = "adminRights"
)

type T_APIHandlerParams struct {
	ResponseWriter   http.ResponseWriter
	Request          *http.Request
	RequestGETParams map[string]string
	RequestURLParams map[string]string
	RequestContent   interface{}
	Session          I_Session
}

type T_RestRedirect struct {
	Result   interface{}
	Location string

	DoRedirect bool
}

type F_APIHandler func(params *T_APIHandlerParams) (interface{}, error)
