package models

import (
	"errors"
)


// retrieves current model implementation and sets its ID to some value
func GetModelAndSetId(modelName string, modelId string) (I_Storable, error) {
	someModel, err := GetModel(modelName)
	if err != nil {
		return nil, err
	}

	storableModel, ok := someModel.(I_Storable)
	if !ok {
		return nil, errors.New("model is not I_Storable capable")
	}

	err = storableModel.SetId(modelId)
	if err != nil {
		return nil, err
	}

	return storableModel, nil
}



// loads model data in current implementation
func LoadModelById(modelName string, modelId string) (I_Storable, error) {

	someModel, err := GetModel(modelName)
	if err != nil {
		return nil, err
	}

	storableModel, ok := someModel.(I_Storable)
	if !ok {
		return nil, errors.New("model is not I_Storable capable")
	}

	err = storableModel.Load(modelId)
	if err != nil {
		return nil, err
	}

	return storableModel, nil
}
