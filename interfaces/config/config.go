package config

import ("errors")

// Interfaces declaration
//-----------------------

type I_Config interface {
	  RegisterItem(Name string, Validator func(interface{}) (interface{}, bool), Default interface{} ) error
	UnregisterItem(Name string) error

	GetValue(Name string) interface{}
	SetValue(Name string, Value interface{}) error

	ListItems() []string

	Load() error
	Save() error
}


// Delegate routines
//------------------

var registeredConfig I_Config

func GetConfig() I_Config { return registeredConfig }

func RegisterConfig(Config I_Config) error {
	if registeredConfig == nil {
		registeredConfig = Config
	} else {
		return errors.New("There is other config already registered")
	}
	return nil
}

func ConfigEmptyValueValidator(val interface{}) (interface{}, bool) { return val, true }


