package models

import (
	"errors"
)

// Package variables
//------------------

var declaredModels = map[string]IModel {}


// Managing routines
//------------------

func RegisterModel(ModelName string, Model IModel) error {
	if _, present := declaredModels[ModelName]; present {
		return errors.New("model with name '" + ModelName + "' already registered")
	} else {
		declaredModels[ModelName] = Model
	}
	return nil
}

func UnRegisterModel(ModelName string) error {
	if _, present := declaredModels[ModelName]; present {
		delete(declaredModels, ModelName)
	} else {
		return errors.New("can't find module with name '" + ModelName + "'")
	}
	return nil
}

func GetModel(ModelName string) (IModel, error) {
	if model, present := declaredModels[ModelName]; present {
		return model.New()
	} else {
		return nil, errors.New("can't find module with name '" + ModelName + "'")
	}
}
