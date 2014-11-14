// Package contain interfaces for API endpoint services.
//
// Currently only "I_RestService" endpoint interface supported.
package api

import (
	"net/http"
)

var (
	SESSION_KEY_ADMIN_RIGHTS = "adminRights" // session key used to flag that user have admin rights
)

// structure to hold API request related information
type T_APIHandlerParams struct {
	ResponseWriter   http.ResponseWriter
	Request          *http.Request
	RequestGETParams map[string]string
	RequestURLParams map[string]string
	RequestContent   interface{}
	Session          I_Session
}

// structure you should return in API handler function if redirect needed
type T_RestRedirect struct {
	Result   interface{}
	Location string

	DoRedirect bool
}

// API handler callback function type
type F_APIHandler func(params *T_APIHandlerParams) (interface{}, error)
