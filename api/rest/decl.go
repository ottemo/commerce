// Package rest is a default implementation of I_RestService declared in "github.com/ottemo/foundation/api" package
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
