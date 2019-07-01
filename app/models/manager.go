package models

import (
	"github.com/ottemo/commerce/env"
)

// RegisterModel registers new model to system
func RegisterModel(ModelName string, Model InterfaceModel) error {
	if _, present := declaredModels[ModelName]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0300eb6b-08b8-497e-afd0-eda0ee358596", "The model with name '"+ModelName+"' has already been registered")
	}
	declaredModels[ModelName] = Model
	return nil
}

// RegisterModel registers new model to system
func RegisterScriptEngine(EngineName string, ScriptEngine InterfaceScriptEngine) error {
	if _, present := declaredScripEngines[EngineName]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "278b6595-29cc-45d8-b599-0e03dae52a46", "Script engine with name '"+EngineName+"' has been already registered")
	}
	declaredScripEngines[EngineName] = ScriptEngine
	return nil
}

// UnRegisterModel removes registered model from system
func UnRegisterModel(ModelName string) error {
	if _, present := declaredModels[ModelName]; present {
		delete(declaredModels, ModelName)
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3d651c0a-e4a4-443f-a6e8-de3a95d89b5c", "Unable to find model to delete with name '"+ModelName+"'")
	}
	return nil
}

// GetModel returns registered in system model
func GetModel(ModelName string) (InterfaceModel, error) {
	if model, present := declaredModels[ModelName]; present {
		return model.New()
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5d49fd0d-1fed-47dc-8e72-2346f1e778c3", "Unable to find model with name '"+ModelName+"'")
}


// GetModel returns registered in system model
func GetScriptEngine(EngineName string) (InterfaceScriptEngine, error) {
	if engine, present := declaredScripEngines[EngineName]; present {
		return engine, nil
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5d49fd0d-1fed-47dc-8e72-2346f1e778c3", "Unable to find script engine with name '"+EngineName+"'")
}


// GetDeclaredModels returns all currently registered in system models
func GetDeclaredModels() map[string]InterfaceModel {
	return declaredModels
}


// GetDeclaredScriptEngines returns all currently registered in system script engines
func GetDeclaredScriptEngines() map[string]InterfaceScriptEngine {
	return declaredScripEngines
}
