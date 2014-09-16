package app

import (
	"errors"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
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
		Description: "application general options",
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
		Description: "application related options",
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

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Store",
		Description: "web store related options",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_NAME,
		Value:       "Ottemo store",
		Type:        "varchar(255)",
		Editor:      "text",
		Options:     nil,
		Label:       "Name",
		Description: "name of your web store",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_EMAIL,
		Value:       "store@ottemo.io",
		Type:        "varchar(255)",
		Editor:      "text",
		Options:     nil,
		Label:       "E-mail",
		Description: "e-mail of your web store",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_COUNTRY,
		Value:       "US",
		Type:        "string",
		Editor:      "select",
		Options:     checkout.COUNTRIES_LIST,
		Label:       "Country",
		Description: "store location country",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_STATE,
		Value:       "",
		Type:        "string",
		Editor:      "select",
		Options:     checkout.STATES_LIST,
		Label:       "State",
		Description: "store location state",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_CITY,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "City",
		Description: "store location city",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_ADDRESSLINE1,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 1",
		Description: "store location address line 1",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_ADDRESSLINE2,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 2",
		Description: "store location address line 2",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_STORE_ZIP,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "zip",
		Description: "store location zip code",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}


	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_MAIL_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Mail",
		Description: "web store mailing options",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_MAIL_SERVER,
		Value:       nil,
		Type:        "varchar(255)",
		Editor:      "text",
		Options:     nil,
		Label:       "Host",
		Description: "web store mailing server",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        CONFIG_PATH_MAIL_PORT,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Port",
		Description: "web store mailing server port",
		Image:       "",
	}, nil)

	if err != nil {
		return err
	}

	return nil
}
