package otto

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/utils"
)

func setupAPI() error {
	service := api.GetRestService()
	service.POST("otto", restOtto)

	return nil
}


// WEB REST API used to execute Otto script
func restOtto(context api.InterfaceApplicationContext) (interface{}, error) {
	// if !api.IsAdminSession(context) {
	// 	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "edabecda-5a46-4745-a8fa-bfd3cb913cb0", "Operation not allowed.")
	// }

	scriptId := ""
	script := ""
	content := context.GetRequestContent()
	if dict, ok := content.(map[string]interface{}); ok {
		if value, present := dict["value"]; present {
			script = utils.InterfaceToString(value)
		}
	} else {
		script = utils.InterfaceToString(content)
	}

	session := context.GetSession()

	if value := session.Get(ConstSessionKey); value != nil {
		scriptId = utils.InterfaceToString(value)
	} else {
		scriptId = utils.MakeUUID()
		session.Set(ConstSessionKey, scriptId)
	}

	vm := engine.GetScriptInstance(scriptId)

	result, err := vm.Execute(script)
	if err != nil {
		return nil, err
	}

	return result, nil
}
