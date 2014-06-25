package json

import (
	"net/http"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/julienschmidt/httprouter"
)

func init() {
	instance := new(DefaultService)

	api.RegisterService(instance)
	env.RegisterOnConfigIniStart(instance.Startup)
}

func (it *DefaultService) Startup() error {

	it.ListenOn = ":9000"
	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("rest.listenOn"); iniValue != "" {
			it.ListenOn = iniValue
		}
	}

	it.Router = httprouter.New()
	it.Router.GET("/",
		func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
			newline := []byte("\n")

			resp.Write([]byte("Ottemo REST API:"))
			resp.Write(newline)

			for path, _ := range it.Handlers {
				resp.Header().Add("Content-Type", "text")
				resp.Write([]byte(path))
				resp.Write(newline)
			}
		})

	it.Handlers = make(map[string]httprouter.Handle)

	api.OnServiceStart()

	return nil
}
