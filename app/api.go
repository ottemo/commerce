package app

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// initializes API for tax
func setupAPI() error {
	var err error = nil

	err = api.GetRestService().RegisterAPI("app", "GET", "login", restLogin)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("app", "POST", "login", restLogin)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("app", "GET", "logout", restLogout)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function login application with root rights
func restLogin(params *api.T_APIHandlerParams) (interface{}, error) {

	var requestLogin string
	var requestPassword string

	if utils.KeysInMapAndNotBlank(params.RequestGETParams, "login", "password") {
		requestLogin = params.RequestGETParams["login"]
		requestPassword = params.RequestGETParams["password"]

	} else {

		reqData, err := api.GetRequestContentAsMap(params)
		if err != nil {
			return nil, err
		}

		if !utils.KeysInMapAndNotBlank(reqData, "login", "password") {
			return nil, errors.New("login and password should be specified")
		}

		requestLogin = utils.InterfaceToString(reqData["login"])
		requestPassword = utils.InterfaceToString(reqData["password"])
	}

	rootLogin := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_STORE_ROOT_LOGIN))
	rootPassword := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_STORE_ROOT_PASSWORD))

	if requestLogin == rootLogin && requestPassword == rootPassword {
		params.Session.Set(api.SESSION_KEY_ADMIN_RIGHTS, true)

		return "ok", nil
	}

	return nil, errors.New("wrong login or password")
}

// WEB REST API function logout application - session data clear
func restLogout(params *api.T_APIHandlerParams) (interface{}, error) {
	err := params.Session.Close()
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
