// Package logger is a default implementation of I_Logger declared in
// "github.com/ottemo/foundation/env" package
package logger

// Package global variables
var (
	baseDirectory     = "./var/log/" // folder location where to store logs
	defaultLogFile    = "system.log" // filename for default log
	defaultErrorsFile = "errors.log" // filename for errors log
)

// I_Logger implementer class
type DefaultLogger struct{}
