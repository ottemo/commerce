package config

import (
	"errors"
	"log"
	"fmt"

	"encoding/json"
	"io/ioutil"

	"github.com/ottemo/platform/tools/module_manager"
	"github.com/ottemo/platform/interfaces/config"
)

func init() {
	module_manager.RegisterModule( new(DefaultConfig) )
}

// structures declaration

type DefaultConfigItem struct {
	Name string
	Validator func(interface{}) (interface{}, bool)
	Default interface{}
	Value interface{}
}

type DefaultConfig struct {
	configValues map[string]*DefaultConfigItem
}



// I_Module interface implementation
//----------------------------------
func (it *DefaultConfig) GetModuleName() string { return "core/Config" }
func (it *DefaultConfig) GetModuleDepends() []string { return make([]string, 0) }

func (it *DefaultConfig) ModuleMakeSysInit() error {
	it.configValues = map[string]*DefaultConfigItem{}
	config.RegisterConfig( it )

	return nil
}
func (it *DefaultConfig) ModuleMakeConfig() error { return nil }
func (it *DefaultConfig) ModuleMakeInit() error {
	it.Load()
	return nil
}
func (it *DefaultConfig) ModuleMakeVerify() error { return nil }
func (it *DefaultConfig) ModuleMakeLoad() error { return nil}
func (it *DefaultConfig) ModuleMakeInstall() error { return nil }
func (it *DefaultConfig) ModuleMakePostInstall() error { return nil }


// I_Config interface implementation
//----------------------------------

func (it *DefaultConfig) RegisterItem(Name string, Validator func(interface{}) (interface{}, bool), Default interface{} ) error {
	if _, present := it.configValues[Name]; present {
		return errors.New("Item [" + Name + "] already registered")
	} else {
		it.configValues[Name] = &DefaultConfigItem{ Name: Name, Validator: Validator, Default: Default, Value: Default }
	}

	return nil
}

func (it *DefaultConfig) UnregisterItem(Name string) error {
	if _, present := it.configValues[Name]; present {
		delete( it.configValues, Name )
	} else {
		return errors.New("Item [" + Name + "] not exists")
	}

	return nil
}

func (it *DefaultConfig) GetValue(Name string) interface{} {
	if configItem, present := it.configValues[Name]; present {
		return configItem.Value
	} else {
		return nil
	}
}

func (it *DefaultConfig) SetValue(Name string, Value interface{}) error {
	if configItem, present := it.configValues[Name]; present {
		if configItem.Validator != nil {

			if newVal, ok := configItem.Validator(Value); ok {
				return errors.New("trying to set invalid value to item [" + Name + "] = " + fmt.Sprintf("%s", Value) )
			} else {
				configItem.Value = newVal
			}

		} else {
			configItem.Value = Value
		}
		return nil
	} else {
		return errors.New("can not find config item [" + Name + "] ")
	}
}

func (it *DefaultConfig) ListItems() []string {
	result := make([]string, len(it.configValues))
	for itemName, _ := range it.configValues {
		result = append( result, itemName )
	}
	return result
}

func (it *DefaultConfig) GetConfigFile() string {
	return "config.conf"
}

func (it *DefaultConfig) Save() error {
	exportMap := map[string]interface{} {}
	for itemName, item := range it.configValues {
		exportMap[itemName] = item.Value
	}

	if data, err := json.Marshal(exportMap); err == nil {
		err := ioutil.WriteFile( it.GetConfigFile(), data, 0666 )
		return err
	} else {
		return err
	}
}

func (it *DefaultConfig) Load() error {

	if data, err := ioutil.ReadFile( it.GetConfigFile() ); err == nil {
		importMap := map[string]interface{} {}
		if err := json.Unmarshal(data, &importMap); err == nil {
			for key, value := range importMap {
				if err := it.SetValue(key, value); err != nil {
					log.Println(err)
				}
			}
		} else {
			return err
		}
	} else {
		return err
	}

	return nil
}
