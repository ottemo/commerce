package app

import (
	"errors"
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
		Label:       "General",
		Description: "Application general options",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_APP_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Application",
		Description: "Application specific options",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STOREFRONT_URL,
		Value:       "http://localhost:8080/",
		Type:        "varchar(255)",
		Editor:      "text",
		Options:     nil,
		Label:       "Storefront host URL",
		Description: "URL application will use to generate storefront resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_DASHBOARD_URL,
		Value:       "http://localhost:9000/",
		Type:        "varchar(255)",
		Editor:      "text",
		Options:     nil,
		Label:       "Dashboard host URL",
		Description: "URL application will use to generate dashboard resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_FOUNDATION_URL,
		Value:       "http://localhost:3000/",
		Type:        "varchar(255)",
		Editor:      "text",
		Options:     nil,
		Label:       "Foundation host URL",
		Description: "URL application will use to generate foundation resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	return nil
}
