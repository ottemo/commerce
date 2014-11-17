package fedex

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_GROUP,
			Value:       nil,
			Type:        env.CONFIG_ITEM_GROUP_TYPE,
			Editor:      "",
			Options:     nil,
			Label:       "FedEx",
			Description: "Federal express shipping method",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_ENABLED,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Enabled",
			Description: "enables/disables shipping method in checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_TITLE,
			Value:       "Federal Express",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Title",
			Description: "shipping method name in checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			if utils.CheckIsBlank(value) {
				return nil, env.ErrorNew("can't be blank")
			} else {
				return value, nil
			}
		})

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_GATEWAY,
			Value:       "https://wsbeta.fedex.com:443/web-services",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Gateway",
			Description: "web services gateway",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_KEY,
			Value:       "",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Account Key",
			Description: "FedEx account key",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_PASSWORD,
			Value:       "",
			Type:        "string",
			Editor:      "password",
			Options:     nil,
			Label:       "Account Password",
			Description: "FedEx account password",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_NUMBER,
			Value:       "",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Account Number",
			Description: "FedEx account number",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_METER,
			Value:       "",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Account Meter",
			Description: "FedEx account meter value",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DEFAULT_WEIGHT,
			Value:       0.1,
			Type:        "decimal",
			Editor:      "decimal",
			Options:     nil,
			Label:       "Default weight",
			Description: "Will be used if product do not have this value (in pounds)",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_ALLOWED_METHODS,
			Value:       "",
			Type:        "string",
			Editor:      "multi_select",
			Options:     SHIPPING_METHODS,
			Label:       "Allowed methods",
			Description: "To customer will be proposed only allowed methods",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DROPOFF,
			Value:       "REGULAR_PICKUP",
			Type:        "string",
			Editor:      "select",
			Options:     SHIPPING_DROPOFF,
			Label:       "Dropoff",
			Description: "dropoff method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			stringValue := utils.InterfaceToString(value)
			if _, present := SHIPPING_DROPOFF[stringValue]; !present {
				return nil, env.ErrorNew("wrong value")
			} else {
				return value, nil
			}
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_PACKAGING,
			Value:       "FEDEX_PAK",
			Type:        "string",
			Editor:      "select",
			Options:     SHIPPING_PACKAGING,
			Label:       "Packing",
			Description: "packing method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			stringValue := utils.InterfaceToString(value)
			if _, present := SHIPPING_PACKAGING[stringValue]; !present {
				return nil, env.ErrorNew("wrong value")
			} else {
				return value, nil
			}
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DEBUG_LOG,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Debug log",
			Description: "enables/disables shipping method debug log",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
