package checkmo

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew("can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigItemGroupType,
		Editor:      "",
		Options:     nil,
		Label:       "Check / Money Order",
		Description: "stub payment method which do nothing until payment",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathEnabled,
		Value:       false,
		Type:        "bool",
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enabled",
		Description: "enables/disables payment method in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTitle,
		Value:       "Check/Money Order",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "Title",
		Description: "payment method name in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			return nil, env.ErrorNew("can't be blank")
		} else {
			return value, nil
		}
	})

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
