package rest

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/utils"
	"io"
)

// GetRequest returns raw request object
func (it *RestApplicationContext) GetRequest() interface{} {
	return it.Request
}

// GetResponse returns raw response object
func (it *RestApplicationContext) GetResponse() interface{} {
	return it.ResponseWriter
}

// GetResponseWriter returns io.Writes for response (for RestApplicationContext it is clone to GetResponse)
func (it *RestApplicationContext) GetResponseWriter() io.Writer {
	return it.ResponseWriter
}

// GetRequestArguments returns arguments required by API function
//   - for REST API it is URI parameters "http://localhost/myfunc/:param1/get/:param2/:param3/"
func (it *RestApplicationContext) GetRequestArguments() map[string]string {
	return it.RequestArguments
}

// GetRequestArgument returns particular argument
func (it *RestApplicationContext) GetRequestArgument(name string) string {
	if value, present := it.RequestArguments[name]; present {
		return value
	}
	return ""
}

// GetRequestParameters returns parameters specified in request
//   - for REST API it is URI options "http://localhost/myfunc/:param1/?option1=1&option2=2"
func (it *RestApplicationContext) GetRequestParameters() map[string]string {
	return it.RequestParameters
}

// GetRequestParameter returns particular request parameter specified in request
func (it *RestApplicationContext) GetRequestParameter(name string) string {
	if value, present := it.RequestParameters[name]; present {
		return value
	}
	return ""
}

// GetRequestContent returns request contents of nil if not specified (HTTP request body)
func (it *RestApplicationContext) GetRequestContent() interface{} {
	return it.RequestContent
}

// GetRequestFiles returns files were attached to request
func (it *RestApplicationContext) GetRequestFiles() map[string]io.Reader {
	return it.RequestFiles
}

// GetRequestFile returns particular file attached to request or nil
func (it *RestApplicationContext) GetRequestFile(name string) io.Reader {
	if file, present := it.RequestFiles[name]; present {
		return file
	}
	return nil
}

// GetRequestSettings returns request related settings
//   - for REST API settings are HTTP headers and COOKIES
func (it *RestApplicationContext) GetRequestSettings() map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range it.Request.Header {
		result[key] = value
	}

	// cookie can overlap header information - it is supposed behaviour
	for _, cookie := range it.Request.Cookies() {
		result[cookie.Name] = cookie.Value
	}

	return result
}

// GetRequestSettings returns particular request related setting or nil
func (it *RestApplicationContext) GetRequestSetting(name string) interface{} {
	if value, err := it.Request.Cookie(name); err == nil && value != nil {
		return value.Value
	}

	if value, present := it.Request.Header[name]; present {
		if len(value) > 1 {
			return value
		}
		return value[0]
	}

	return nil
}

// GetRequestContentType returns MIME type of request content
func (it *RestApplicationContext) GetRequestContentType() string {
	return it.Request.Header.Get("Content-Type")
}

// GetResponseContentType returns MIME type of supposed response content
func (it *RestApplicationContext) GetResponseContentType() string {
	return it.ResponseWriter.Header().Get("Content-Type")
}

// SetResponseContentType changes response content type, returns error if not possible
func (it *RestApplicationContext) SetResponseContentType(mimeType string) error {
	it.ResponseWriter.Header().Set("Content-Type", mimeType)
	return nil
}

func (it *RestApplicationContext) GetResponseSetting(name string) interface{} {
	return it.ResponseWriter.Header().Get(name)
}

// SetResponseSetting specifies response setting (for REST API it just sets additional header)
func (it *RestApplicationContext) SetResponseSetting(name string, value interface{}) error {
	it.ResponseWriter.Header().Set(name, utils.InterfaceToString(value))
	return nil
}

// GetResponseResult returns result going to be written to response writer
func (it *RestApplicationContext) GetResponseResult() interface{} {
	return it.Result
}

// SetResponseResult changes result going to be written to response writer
func (it *RestApplicationContext) SetResponseResult(value interface{}) error {
	it.Result = value
	return nil
}

// GetContextValues returns current context related values map
func (it *RestApplicationContext) GetContextValues() map[string]interface{} {
	return it.ContextValues
}

// GetContextValue returns particular context related value or nil if not set
func (it *RestApplicationContext) GetContextValue(key string) interface{} {
	if value, present := it.ContextValues[key]; present {
		return value
	}
	return nil
}

// SetContextValue stores specified value in current context
func (it *RestApplicationContext) SetContextValue(key string, value interface{}) {
	it.ContextValues[key] = value
}

// SetSession assigns given session to current context
func (it *RestApplicationContext) SetSession(session api.InterfaceSession) error {
	it.Session = session
	return nil
}

// GetSession returns session assigned to current context or nil if nothing was assigned
func (it *RestApplicationContext) GetSession() api.InterfaceSession {
	return it.Session
}
