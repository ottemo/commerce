// Package "logger" is a default implementation for "I_Logger" interface.
package logger

var (
	baseDirectory     = "./var/log/" // folder location where to store logs
	defaultLogFile    = "system.log" // filename for default log
	defaultErrorsFile = "errors.log" // filename for errors log
)

// I_Logger implementer class
type DefaultLogger struct{}
