package rest

import (
	"net/http"

	"github.com/ottemo/foundation/config"
	"github.com/ottemo/foundation/api"

	"github.com/julienschmidt/httprouter"
)

func init() {
	instance := new(DefaultRestService)

	api.RegisterRestService(instance)
	config.RegisterOnConfigIniStart( instance.startup )
}

func (it *DefaultRestService) startup() error {

	it.ListenOn = ":3000"
	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("rest.listenOn"); iniValue != "" {
			it.ListenOn = iniValue
		}
	}

	it.Router = httprouter.New()
	it.Router.GET("/",
		func( resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
			newline := []byte( "\n" )

			resp.Write( []byte( "Ottemo REST API:" ) )
			resp.Write( newline )

			for path, _ := range it.Handlers {
				resp.Header().Add("Content-Type", "text")
				resp.Write( []byte( path ) )
				resp.Write( newline )
			}
		})

	it.Handlers = make( map[string]httprouter.Handle )

	api.OnRestServiceStart()

	return nil
}
