package config

import (
	"errors"
)

var registeredConfig I_Config
var registeredIniConfig I_IniConfig
var callbacksOnConfigStart = []func() error {}
var callbacksOnConfigIniStart = []func() error {}

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

func GetConfig() I_Config { return registeredConfig }
func GetIniConfig() I_IniConfig { return registeredIniConfig }



func ConfigEmptyValueValidator(val interface{}) (interface{}, bool) { return val, true }
