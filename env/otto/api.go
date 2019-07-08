package otto

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/env"
)

func setupAPI() error {
	service := api.GetRestService()
	service.POST("otto", restOtto)

	return nil
}


// WEB REST API used to execute Otto script
func restOtto(context api.InterfaceApplicationContext) (interface{}, error) {
	if !api.IsAdminSession(context) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "edabecda-5a46-4745-a8fa-bfd3cb913cb0", "Operation not allowed.")
	}



	return nil, nil
}
