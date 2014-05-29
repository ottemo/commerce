package web_server

import (
	"errors"
)

// Interfaces declaration
//-----------------------

type I_WebServer interface {
	GetName() string
	Run() error

	RegisterController(HTTPType string, Path string, CallbackFunc interface{}) error  // TODO: bad planned interaction ( CallbackFunc - impossible to imagine )
}

// Delegate routines
//------------------

var registeredWebServer I_WebServer

func GetWebServer() I_WebServer {
	return registeredWebServer
}

func RegisterWebServer(WebServer I_WebServer) error {
	if registeredWebServer == nil {
		registeredWebServer = WebServer
	} else {
		return errors.New("The web server '" + registeredWebServer.GetName() + "' already registered")
	}
	return nil
}
