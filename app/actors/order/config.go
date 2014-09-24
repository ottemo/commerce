package order

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupConfig() error {
	config := env.GetConfig()

	config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_LAST_INCREMENT_ID,
		Value:       0,
		Type:        "int",
		Editor:      "integer",
		Options:     "",
		Label:       "Last Order Increment ID: ",
		Description: "Do not change this value unless you know what you doing",
		Image:       "",
	},
		func(value interface{}) (interface{}, error) {
			return utils.InterfaceToInt(value), nil
		})

	lastIncrementId = utils.InterfaceToInt(config.GetValue(CONFIG_PATH_LAST_INCREMENT_ID))

	return nil
}
