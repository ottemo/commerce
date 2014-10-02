package env

import (
	"errors"
)

var (
	registeredConfig    I_Config
	registeredIniConfig I_IniConfig

	registeredLogger   I_Logger
	registeredErrorBus I_ErrorBus
	registeredEventBus I_EventBus

	callbacksOnConfigStart    = []func() error{}
	callbacksOnConfigIniStart = []func() error{}
)

// callbacks
//----------

func RegisterOnConfigStart(callback func() error) {
	callbacksOnConfigStart = append(callbacksOnConfigStart, callback)
}

func OnConfigStart() error {
	for _, callback := range callbacksOnConfigStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

func RegisterOnConfigIniStart(callback func() error) {
	callbacksOnConfigIniStart = append(callbacksOnConfigIniStart, callback)
}

func OnConfigIniStart() error {
	for _, callback := range callbacksOnConfigIniStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// objects
//---------

func RegisterIniConfig(IniConfig I_IniConfig) error {
	if registeredIniConfig == nil {
		registeredIniConfig = IniConfig
	} else {
		return errors.New("There is other ini config already registered")
	}
	return nil
}

func RegisterConfig(Config I_Config) error {
	if registeredConfig == nil {
		registeredConfig = Config
	} else {
		return errors.New("There is other config already registered")
	}
	return nil
}

func RegisterLogger(logger I_Logger) error {
	if registeredLogger == nil {
		registeredLogger = logger
	} else {
		return errors.New("Logger already registered")
	}
	return nil
}

func RegisterEventBus(eventBus I_EventBus) error {
	if registeredEventBus == nil {
		registeredEventBus = eventBus
	} else {
		return errors.New("Event bus already registered")
	}
	return nil
}

func RegisterErrorBus(errorBus I_ErrorBus) error {
	if registeredErrorBus == nil {
		registeredErrorBus = errorBus
	} else {
		return errors.New("Error bus already registered")
	}
	return nil
}

func GetConfig() I_Config       { return registeredConfig }
func GetIniConfig() I_IniConfig { return registeredIniConfig }
func GetLogger() I_Logger       { return registeredLogger }
func GetErrorBus() I_ErrorBus   { return registeredErrorBus }
func GetEventBus() I_EventBus   { return registeredEventBus }

func ConfigEmptyValueValidator(val interface{}) (interface{}, bool) { return val, true }
