package negroni

import (
	"net/http"
	base_negroni "github.com/codegangsta/negroni"
)

type HTTPHandler func(resp http.ResponseWriter, req *http.Request)
type JSONHandler func(req *http.Request) map[string]interface{}

type NegroniRestService struct {
	Negroni  *base_negroni.Negroni
	Mux		 *http.ServeMux
	ListenOn string

	Handlers map[string]HTTPHandler
}


