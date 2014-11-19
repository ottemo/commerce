// Package logger is a default implementation of InterfaceLogger declared in
// "github.com/ottemo/foundation/env" package
package logger

// Package global variables
var (
	baseDirectory     = "./var/log/" // folder location where to store logs
	defaultLogFile    = "system.log" // filename for default log
	defaultErrorsFile = "errors.log" // filename for errors log
)

// InterfaceLogger implementer class
type DefaultLogger struct{}
