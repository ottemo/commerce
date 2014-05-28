package config

import ("errors")

// variables and declarations
//----------------------------------

type I_Config interface {
	  RegisterItem(Name string, TypeName string, Validator func(interface{}) bool, Default interface{} ) error
	UnregisterItem(Name string) error

	GetValue(Name string) interface{}
	SetValue(Name string, Value interface{}) error

	ListItems() []string

	Load() error
	Save() error
}

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

func ConfigEmptyValueValidator(interface{}) bool { return true }


