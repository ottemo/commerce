package negroni

import (
	"net/http"

	"github.com/ottemo/foundation/config"
	"github.com/ottemo/foundation/rest_service"

	base_negroni "github.com/codegangsta/negroni"
)

func init() {
	instance := new(NegroniRestService)

	rest_service.RegisterRestService(instance)
	config.RegisterOnConfigIniStart( instance.Startup )
}

func (it *NegroniRestService) Startup() error {
	it.Mux = http.NewServeMux()

	it.ListenOn = ":9000"
	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("negroni.port"); iniValue != "" {
			it.ListenOn = iniValue
		}
	}

	it.Mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		newline := []byte( "\n" )
		for path, _ := range it.Handlers {
			resp.Header().Add("Content-Type", "text")
			resp.Write( []byte( path ) )
			resp.Write( newline )
		}
	})

	it.Negroni = base_negroni.Classic()
	it.Negroni.UseHandler(it.Mux)

	it.Handlers = make( map[string]HTTPHandler )

	rest_service.OnRestServiceStart()

	return nil
}
