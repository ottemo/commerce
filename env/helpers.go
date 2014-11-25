package env

import (
	"errors"
)

// ConfigGetValue returns config value or nil if not present
func ConfigGetValue(Path string) interface{} {
	if config := GetConfig(); config != nil {
		return config.GetValue(Path)
	}

	return nil
}

// IniValue returns value from ini file or "" if not present
func IniValue(Path string) string {
	if iniConfig := GetIniConfig(); iniConfig != nil {
		return iniConfig.GetValue(Path, "")
	}
	return ""
}

// Log logs general purpose message
func Log(storage string, prefix string, message string) {
	if logger := GetLogger(); logger != nil {
		logger.Log(storage, prefix, message)
	}
}

// LogError logs an error message
func LogError(err error) {
	if logger := GetLogger(); logger != nil {
		logger.LogError(err)
	}
}

// LogMessage is a Log function short form for info messages in default storage
func LogMessage(message string) {
	if logger := GetLogger(); logger != nil {
		logger.LogMessage(message)
	}
}

// ErrorLevel returns error level of given error
func ErrorLevel(err error) int {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorLevel(err)
	}
	return 0
}

// ErrorCode returns error code of given error
func ErrorCode(err error) string {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorCode(err)
	}
	return ""
}

// ErrorMessage returns message of given error
func ErrorMessage(err error) string {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorMessage(err)
	}
	return err.Error()
}

// ErrorRegisterListener registers listener for error bus
func ErrorRegisterListener(listener FuncErrorListener) {
	if errorBus := GetErrorBus(); errorBus != nil {
		errorBus.RegisterListener(listener)
	}
}

// ErrorDispatch handles error, returns new one you should pass next
func ErrorDispatch(err error) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.Dispatch(err)
	}
	return err
}

// ErrorNew creates new error and dispatches it
func ErrorNew(message string) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.New(message)
	}
	return errors.New(message)
}

// EventRegisterListener registers listener for event bus
func EventRegisterListener(event string, listener FuncEventListener) {
	if eventBus := GetEventBus(); eventBus != nil {
		eventBus.RegisterListener(event, listener)
	}
}

// Event emits new event for registered listeners
func Event(event string, args map[string]interface{}) {
	if eventBus := GetEventBus(); eventBus != nil {
		eventBus.New(event, args)
	}
}
