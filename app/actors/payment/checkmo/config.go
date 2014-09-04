package checkmo

import (
	"errors"

	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/env"
)

// setup configuration values
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return errors.New("can't obtain config")
	}

	err := config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Check / Money Order",
		Description: "stub payment method which do nothing until payment",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_ENABLED,
		Value:       false,
		Type:        "bool",
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enabled",
		Description: "enables/disables payment method in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_TITLE,
		Value:       "Check/Money Order",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "Title",
		Description: "payment method name in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			return nil, errors.New("can't be blank")
		} else {
			return value, nil
		}
	})

	if err != nil {
		return err
	}

	return nil
}
