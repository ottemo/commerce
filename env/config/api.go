package config

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
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
