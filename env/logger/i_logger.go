package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/ottemo/foundation/env"
)

// general logging function
func (it *DefaultLogger) Log(storage string, prefix string, msg string) {
	message := time.Now().Format(time.RFC3339) + ": " + msg + "\n"

	logFile, err := os.OpenFile(baseDirectory+storage, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(message)
		return
	}

	logFile.Write([]byte(message))

	logFile.Close()
}

// makes error log
func (it *DefaultLogger) LogError(err error) {
	if ottemoErr, ok := err.(env.InterfaceOttemoError); ok {
		it.Log(defaultErrorsFile, env.ConstLogPrefixError, ottemoErr.ErrorFull())
	} else {
		it.Log(defaultErrorsFile, env.ConstLogPrefixError, err.Error())
	}
}

// log message to separate file
func (it *DefaultLogger) LogToStorage(storage string, msg string) {
	it.Log(storage, env.ConstLogPrefixInfo, msg)
}

// log message with prefix specification
func (it *DefaultLogger) LogWithPrefix(prefix string, msg string) {
	it.Log(defaultLogFile, prefix, msg)
}

// simplified logging function
func (it *DefaultLogger) LogMessage(msg string) {
	it.Log(defaultLogFile, env.ConstLogPrefixInfo, msg)
}
