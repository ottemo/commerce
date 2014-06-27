package default_rest_service

import (
	"net/http"

	"github.com/ottemo/foundation/config"
	"github.com/ottemo/foundation/rest_service"

	"github.com/julienschmidt/httprouter"
)

func init() {
	instance := new(DefaultRestService)

	rest_service.RegisterRestService(instance)
	config.RegisterOnConfigIniStart( instance.Startup )
}

func (it *DefaultRestService) Startup() error {

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

	rest_service.OnRestServiceStart()

	return nil
}
