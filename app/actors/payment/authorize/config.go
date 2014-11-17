package authorize

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_GROUP,
			Value:       nil,
			Type:        env.CONFIG_ITEM_GROUP_TYPE,
			Editor:      "",
			Options:     nil,
			Label:       "Authorize.Net (Direct Post)",
			Description: "see https://developer.authorize.net/api/dpm/",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_ENABLED,
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

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_TITLE,
			Value:       "Authorize.Net (Direct Post)",
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

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_ACTION,
			Value:       DPM_ACTION_AUTHORIZE_ONLY,
			Type:        "string",
			Editor:      "select",
			Options:     map[string]string{DPM_ACTION_AUTHORIZE_ONLY: "Authorize Only", DPM_ACTION_AUTHORIZE_AND_CAPTURE: "Authorize & Capture"},
			Label:       "Action",
			Description: "specifies action on checkout submit",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			stringValue := utils.InterfaceToString(value)
			if !utils.IsAmongStr(stringValue, DPM_ACTION_AUTHORIZE_ONLY, DPM_ACTION_AUTHORIZE_AND_CAPTURE) {
				return nil, env.ErrorNew("should be " + DPM_ACTION_AUTHORIZE_ONLY + " or " + DPM_ACTION_AUTHORIZE_AND_CAPTURE)
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_LOGIN,
			Value:       "",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "API Login ID",
			Description: "account login id",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_KEY,
			Value:       "",
			Type:        "string",
			Editor:      "password",
			Options:     nil,
			Label:       "Transaction Key",
			Description: "account transaction key",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_GATEWAY,
			Value:       "https://test.authorize.net/gateway/transact.dll",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Gateway",
			Description: "payment method gateway",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_TEST,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Test Mode",
			Description: "specifies test mode for payment method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_DEBUG,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Debug",
			Description: "specifies to write log on that payment method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DPM_CHECKOUT,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Custom checkout page",
			Description: "use the custom relay page on checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return err
		}
	}

	return nil
}
