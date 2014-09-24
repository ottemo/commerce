package config

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("config", "GET", "groups", restConfigGroups)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "info/:path", restConfigInfo)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "list", restConfigList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "get/:path", restConfigGet)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("config", "POST", "set/:path", restConfigSet)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("config", "POST", "register", restConfigRegister)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("config", "DELETE", "unregister/:path", restConfigUnRegister)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "reload", restConfigReload)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API to get value information about config items with type [CONFIG_ITEM_GROUP_TYPE]
func restConfigGroups(params *api.T_APIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.GetGroupItems(), nil
}

// WEB REST API to get value information about config items with type [CONFIG_ITEM_GROUP_TYPE]
func restConfigList(params *api.T_APIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.ListPathes(), nil
}

// WEB REST API to get value information about item(s) matching path
func restConfigInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.GetItemsInfo(params.RequestURLParams["path"]), nil
}

// WEB REST API used to get value of particular item in config
//   - path should be without any wildcard
func restConfigGet(params *api.T_APIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.GetValue(params.RequestURLParams["path"]), nil
}

// WEB REST API used to set value of particular item in config
//   - path should be without any wildcard
func restConfigSet(params *api.T_APIHandlerParams) (interface{}, error) {

	config := env.GetConfig()

	var setValue interface{} = nil

	setValue = params.RequestContent
	configPath := params.RequestURLParams["path"]

	content, err := api.GetRequestContentAsMap(params)
	if err == nil {
		if contentValue, present := content["value"]; present {
			setValue = contentValue
		}
	}

	err = config.SetValue(configPath, setValue)

	return config.GetValue(configPath), err
}

// WEB REST API used to add new config Item to a config system
func restConfigRegister(params *api.T_APIHandlerParams) (interface{}, error) {
	inputData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	config := env.GetConfig()

	configItem := env.T_ConfigItem{
		Path:  utils.InterfaceToString(utils.GetFirstMapValue(inputData, "path", "Path")),
		Value: utils.GetFirstMapValue(inputData, "value"),

		Type: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "type", "Type")),

		Editor:  utils.InterfaceToString(utils.GetFirstMapValue(inputData, "editor", "Editor")),
		Options: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "options", "Options")),

		Label:       utils.InterfaceToString(utils.GetFirstMapValue(inputData, "label", "Label")),
		Description: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "description", "Description")),

		Image: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "image", "Image")),
	}

	config.RegisterItem(configItem, nil)

	return configItem, nil
}

// WEB REST API used to remove config item from system
func restConfigUnRegister(params *api.T_APIHandlerParams) (interface{}, error) {
	config := env.GetConfig()

	err := config.UnregisterItem(params.RequestURLParams["path"])
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API used to re-load config from DB
func restConfigReload(params *api.T_APIHandlerParams) (interface{}, error) {
	config := env.GetConfig()

	err := config.Reload()
	if err != nil {
		return nil, err
	}

	return "ok", nil
}
