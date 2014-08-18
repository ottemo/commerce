package flat

import (
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/env"
)

// module entry point before app start
func init() {
	instance := new(FlatRateShipping)

	checkout.RegisterShippingMethod(instance)

	env.RegisterOnConfigStart(setupConfig)
}

func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_GROUP,
			Value:       nil,
			Type:        env.CONFIG_ITEM_GROUP_TYPE,
			Editor:      "",
			Options:     nil,
			Label:       "Flat Rate",
			Description: "static amount stipping method",
			Image:       "",
		}, nil)

		if err != nil {
			return err
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_ENABLED,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Enabled",
			Description: "enables/disables shipping method for storefront",
			Image:       "",
		}, nil)

		if err != nil {
			return err
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_AMOUNT,
			Value:       10,
			Type:        "int",
			Editor:      "money",
			Options:     nil,
			Label:       "Amount",
			Description: "price of shipping",
			Image:       "",
		}, nil)

		if err != nil {
			return err
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_NAME,
			Value:       "Flat Rate",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Name",
			Description: "shipping name displayed in checkout",
			Image:       "",
		}, nil)

		if err != nil {
			return err
		}

		config.RegisterItem(env.T_ConfigItem{
			Path:        CONFIG_PATH_DAYS,
			Value:       0,
			Type:        "int",
			Editor:      "integer",
			Options:     nil,
			Label:       "Ship days",
			Description: "amount of days for shipping",
			Image:       "",
		}, nil)

		if err != nil {
			return err
		}
	}

	return nil
}
