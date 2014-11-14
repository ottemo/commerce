// Package "rest" is a default implementation for "I_RestService" interface.
package rest

import (
	"github.com/julienschmidt/httprouter"
)

// I_RestService implementer class
type DefaultRestService struct {
	ListenOn string
	Router   *httprouter.Router
	Handlers map[string]httprouter.Handle
}
