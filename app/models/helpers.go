package models

import (
	"github.com/ottemo/foundation/env"
)

// retrieves current model implementation and sets its ID to some value
func GetModelAndSetId(modelName string, modelId string) (I_Storable, error) {
	someModel, err := GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storableModel, ok := someModel.(I_Storable)
	if !ok {
		return nil, env.ErrorNew("model is not I_Storable capable")
	}

	err = storableModel.SetId(modelId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return storableModel, nil
}

// loads model data in current implementation
func LoadModelById(modelName string, modelId string) (I_Storable, error) {

	someModel, err := GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storableModel, ok := someModel.(I_Storable)
	if !ok {
		return nil, env.ErrorNew("model is not I_Storable capable")
	}

	err = storableModel.Load(modelId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return storableModel, nil
}
