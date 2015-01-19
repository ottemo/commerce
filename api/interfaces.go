package api

import (
	"io"
	"net/http"
)

// InterfaceSessionService is an interface to access session managing service
type InterfaceSessionService interface {
	GetName() string

	New() (InterfaceSession, error)
	Get(sessionID string) (InterfaceSession, error)
}

// InterfaceSession is an interface represents private storage for particular API request
type InterfaceSession interface {
	GetID() string

	Get(key string) interface{}
	Set(key string, value interface{})

	SetModified()

	Close() error

	Load(id string) error
	Save() error
}

// InterfaceRestService is an interface to interact with RESTFul API service
type InterfaceRestService interface {
	GetName() string

	Run() error
	RegisterAPI(service string, method string, uri string, handler FuncAPIHandler) error

	http.Handler
}

// InterfaceApplicationContextSupport is an interface to assign/get application context to object
type InterfaceApplicationContextSupport interface {
	GetApplicationContext() InterfaceApplicationContext
	SetApplicationContext(context InterfaceApplicationContext) error
}

// InterfaceApplicationContext is an interface representing context where current execution happens
type InterfaceApplicationContext interface {
	GetRequest() interface{}
	GetResponse() interface{}

	GetSession() InterfaceSession

	GetContextValues() map[string]interface{}
	GetContextValue(key string) interface{}
	SetContextValue(key string, value interface{})

	GetResponseWriter() io.Writer

	GetRequestArguments() map[string]string
	GetRequestArgument(name string) string
	GetRequestParameters() map[string]string
	GetRequestParameter(name string) string
	GetRequestFiles() map[string]io.Reader
	GetRequestFile(name string) io.Reader
	GetRequestSettings() map[string]interface{}
	GetRequestSetting(name string) interface{}
	GetRequestContent() interface{}
	GetRequestContentType() string

	GetResponseContentType() string
	SetResponseContentType(mimeType string) error
	GetResponseSetting(name string) interface{}
	SetResponseSetting(name string, value interface{}) error
	GetResponseResult() interface{}
	SetResponseResult(value interface{}) error
}
