package rest

import (
	"github.com/julienschmidt/httprouter"
)

type DefaultRestService struct {
	ListenOn string
	Router   *httprouter.Router
	Handlers map[string]httprouter.Handle
}


