package logger

var (
	baseDirectory     = "./var/log/"
	defaultLogFile    = "system.log"
	defaultErrorsFile = "errors.log"
)

type DefaultLogger struct{}
