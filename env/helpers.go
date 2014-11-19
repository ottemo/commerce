package env

import (
	"errors"
)

// returns config value or nil if not present
func ConfigGetValue(Path string) interface{} {
	if config := GetConfig(); config != nil {
		return config.GetValue(Path)
	}

	return nil
}

// returns value from ini file or "" if not present
func IniValue(Path string) string {
	if iniConfig := GetIniConfig(); iniConfig != nil {
		return iniConfig.GetValue(Path, "")
	}
	return ""
}

// logs general purpose message
func Log(storage string, prefix string, message string) {
	if logger := GetLogger(); logger != nil {
		logger.Log(storage, prefix, message)
	}
}

// log for error
func LogError(err error) {
	if logger := GetLogger(); logger != nil {
		logger.LogError(err)
	}
}

// log short form for info message
func LogMessage(message string) {
	if logger := GetLogger(); logger != nil {
		logger.LogMessage(message)
	}
}

// returns error level of given error
func ErrorLevel(err error) int {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorLevel(err)
	}
	return 0
}

// returns error code of given error
func ErrorCode(err error) string {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorCode(err)
	}
	return ""
}

// returns message of given error
func ErrorMessage(err error) string {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorMessage(err)
	}
	return err.Error()
}

// registers listener for error bus
func ErrorRegisterListener(listener FuncErrorListener) {
	if errorBus := GetErrorBus(); errorBus != nil {
		errorBus.RegisterListener(listener)
	}
}

// handles error, returns new one you should pass next
func ErrorDispatch(err error) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.Dispatch(err)
	}
	return err
}

// creates new error and dispatches it
func ErrorNew(message string) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.New(message)
	}
	return errors.New(message)
}

// registers listener for event bus
func EventRegisterListener(event string, listener FuncEventListener) {
	if eventBus := GetEventBus(); eventBus != nil {
		eventBus.RegisterListener(event, listener)
	}
}

// emits new event for registered listeners
func Event(event string, args map[string]interface{}) {
	if eventBus := GetEventBus(); eventBus != nil {
		eventBus.New(event, args)
	}
}
