// Package rest is a default implementation of InterfaceRestService
// declared in "github.com/ottemo/foundation/api" package
package rest

import (
	"github.com/julienschmidt/httprouter"
)

// DefaultRestService is a default implementer of InterfaceRestService declared in
// "github.com/ottemo/foundation/api" package
type DefaultRestService struct {
	ListenOn string
	Router   *httprouter.Router
	Handlers map[string]httprouter.Handle
}
