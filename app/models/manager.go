package models

import (
	"github.com/ottemo/foundation/env"
)

func RegisterModel(ModelName string, Model I_Model) error {
	if _, present := declaredModels[ModelName]; present {
		return env.ErrorNew("model with name '" + ModelName + "' already registered")
	} else {
		declaredModels[ModelName] = Model
	}
	return nil
}

func UnRegisterModel(ModelName string) error {
	if _, present := declaredModels[ModelName]; present {
		delete(declaredModels, ModelName)
	} else {
		return env.ErrorNew("can't find module with name '" + ModelName + "'")
	}
	return nil
}

func GetModel(ModelName string) (I_Model, error) {
	if model, present := declaredModels[ModelName]; present {
		return model.New()
	} else {
		return nil, env.ErrorNew("can't find module with name '" + ModelName + "'")
	}
}
