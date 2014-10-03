package paypal

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setup configuration values
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew("can't obtain config")
	}

	err := config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "PayPal (Express)",
		Description: "see http://paypal.com",
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
		Value:       "PayPal",
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
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_NVP,
		Value:       "https://api-3t.sandbox.paypal.com/nvp",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "NVP Gateway",
		Description: "URL to send Name-Value Pair request",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_GATEWAY,
		Value:       "https://www.sandbox.paypal.com/webscr?cmd=_express-checkout",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "Gateway",
		Description: "PayPal gateway",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_USER,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "User",
		Description: "PayPal username",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_PASS,
		Value:       "",
		Type:        "string",
		Editor:      "password",
		Options:     nil,
		Label:       "Password",
		Description: "PayPal password",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_SIGNATURE,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "Signature",
		Description: "PayPal signature",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_ACTION,
		Value:       "",
		Type:        "string",
		Editor:      "select",
		Options:     map[string]string{PAYMENT_ACTION_SALE: "Sale", PAYMENT_ACTION_AUTHORIZATION: "Authorization"},
		Label:       "Signature",
		Description: "PayPal signature",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	return nil
}
