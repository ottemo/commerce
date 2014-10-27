package impex

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

type ImpexImportCmdInsert struct {
	Model models.I_Model
}

type ImpexImportCmdUpdate struct{}
type ImpexImportCmdDelete struct{}
type ImpexImportCmdStore struct{}

func (it *ImpexImportCmdInsert) Init(args []string, exchange map[string]interface{}) error {
	if len(args) > 1 {
		modelName := args[1]

		cmdModel, err := models.GetModel(modelName)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		if _, ok := cmdModel.(models.I_Object); !ok {
			return env.ErrorNew("model '" + modelName + "' is not I_Object")
		}
		if _, ok := cmdModel.(models.I_Storable); !ok {
			return env.ErrorNew("model '" + modelName + "' is not I_Storable")
		}

		it.Model = cmdModel
	} else {
		return env.ErrorNew("model insert into was not specified")
	}
	return nil
}

func (it *ImpexImportCmdInsert) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	cmdModel, err := it.Model.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	if modelObject, ok := cmdModel.(models.I_Object); ok {
		for attribute, value := range itemData {
			modelObject.Set(attribute, value)
		}

		if modelStorable, ok := cmdModel.(models.I_Storable); ok {
			err := modelStorable.Save()
			if err != nil {
				return nil, err
			}
		}
	}
	return cmdModel, nil
}

func (it *ImpexImportCmdUpdate) Init(args []string, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdUpdate) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdDelete) Init(args []string, exchange map[string]interface{}) error {
	return env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdDelete) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdStore) Init(args []string, exchange map[string]interface{}) error {
	return nil //env.ErrorNew("not implemented")
}

func (it *ImpexImportCmdStore) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	return nil, nil //nil, env.ErrorNew("not implemented")
}
