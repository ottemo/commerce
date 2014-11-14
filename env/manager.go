package env

import (
	"errors"
)

var (
	// variables to hold currently registered services
	registeredConfig    I_Config
	registeredIniConfig I_IniConfig
	registeredLogger    I_Logger
	registeredErrorBus  I_ErrorBus
	registeredEventBus  I_EventBus

	// variables to hold callback functions on configuration services startup
	callbacksOnConfigStart    = []func() error{}
	callbacksOnConfigIniStart = []func() error{}
)

// registers new callback on configuration service start
func RegisterOnConfigStart(callback func() error) {
	callbacksOnConfigStart = append(callbacksOnConfigStart, callback)
}

// registers new callback on ini configuration service start
func RegisterOnConfigIniStart(callback func() error) {
	callbacksOnConfigIniStart = append(callbacksOnConfigIniStart, callback)
}

// fires config service start event (callback handling)
func OnConfigStart() error {
	for _, callback := range callbacksOnConfigStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// fires ini config service start event (callback handling)
func OnConfigIniStart() error {
	for _, callback := range callbacksOnConfigIniStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// registers ini config service in the system
//   - will cause error if there are couple candidates for that role
func RegisterIniConfig(IniConfig I_IniConfig) error {
	if registeredIniConfig == nil {
		registeredIniConfig = IniConfig
	} else {
		return errors.New("There is other ini config already registered")
	}
	return nil
}

// registers config service in the system
//   - will cause error if there are couple candidates for that role
func RegisterConfig(Config I_Config) error {
	if registeredConfig == nil {
		registeredConfig = Config
	} else {
		return errors.New("There is other config already registered")
	}
	return nil
}

// registers logging service in the system
//   - will cause error if there are couple candidates for that role
func RegisterLogger(logger I_Logger) error {
	if registeredLogger == nil {
		registeredLogger = logger
	} else {
		return errors.New("Logger already registered")
	}
	return nil
}

// registers event processor in the system
//   - will cause error if there are couple candidates for that role
func RegisterEventBus(eventBus I_EventBus) error {
	if registeredEventBus == nil {
		registeredEventBus = eventBus
	} else {
		return errors.New("Event bus already registered")
	}
	return nil
}

// registers error processor in the system
//   - will cause error if there are couple candidates for that role
func RegisterErrorBus(errorBus I_ErrorBus) error {
	if registeredErrorBus == nil {
		registeredErrorBus = errorBus
	} else {
		return errors.New("Error bus already registered")
	}
	return nil
}

// returns currently used config service implementation
func GetConfig() I_Config {
	return registeredConfig
}

// returns currently used ini config service implementation
func GetIniConfig() I_IniConfig {
	return registeredIniConfig
}

// returns currently used logging service implementation
func GetLogger() I_Logger {
	return registeredLogger
}

// returns currently used error processor implementation
func GetErrorBus() I_ErrorBus {
	return registeredErrorBus
}

// returns currently used event processor implementation
func GetEventBus() I_EventBus {
	return registeredEventBus
}

// validator function to accept any value
func ConfigEmptyValueValidator(val interface{}) (interface{}, bool) {
	return val, true
}
