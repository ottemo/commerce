package config

import (
	"errors"
	"fmt"

	"github.com/ottemo/foundation/env"
)

func (it *DefaultConfig) RegisterItem(Name string, Validator func(interface{}) (interface{}, bool), Default interface{}) error {
	if _, present := it.configValues[Name]; present {
		return errors.New("Item [" + Name + "] already registered")
	} else {
		it.configValues[Name] = &DefaultConfigItem{Name: Name, Validator: Validator, Default: Default, Value: Default}
	}

	return nil
}

func (it *DefaultConfig) UnregisterItem(Name string) error {
	if _, present := it.configValues[Name]; present {
		delete(it.configValues, Name)
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
				return errors.New("trying to set invalid value to item [" + Name + "] = " + fmt.Sprintf("%s", Value))
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
		result = append(result, itemName)
	}
	return result
}

func (it *DefaultConfig) Save() error {
	return nil
}

func (it *DefaultConfig) Load() error {

	env.OnConfigStart()

	return nil
}
