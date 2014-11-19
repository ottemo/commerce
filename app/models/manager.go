package models

import (
	"github.com/ottemo/foundation/env"
)

// RegisterModel registers new model to system
func RegisterModel(ModelName string, Model InterfaceModel) error {
	if _, present := declaredModels[ModelName]; present {
		return env.ErrorNew("model with name '" + ModelName + "' already registered")
	}
	declaredModels[ModelName] = Model

	return nil
}

// UnRegisterModel removes registered model from system
func UnRegisterModel(ModelName string) error {
	if _, present := declaredModels[ModelName]; present {
		delete(declaredModels, ModelName)
	} else {
		return env.ErrorNew("can't find module with name '" + ModelName + "'")
	}
	return nil
}

// GetModel returns registered in system model
func GetModel(ModelName string) (InterfaceModel, error) {
	if model, present := declaredModels[ModelName]; present {
		return model.New()
	}
	return nil, env.ErrorNew("can't find module with name '" + ModelName + "'")
}

// GetDeclaredModels returns all currently registered in system models
func GetDeclaredModels() map[string]InterfaceModel {
	return declaredModels
}
