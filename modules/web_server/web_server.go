package web_server

import (
	"errors"
	"strconv"
	"log"

	"github.com/ottemo/platform/interfaces/web_server"
	"github.com/ottemo/platform/interfaces/config"

	"github.com/ottemo/platform/tools/module_manager"

	"net/http"
	"github.com/julienschmidt/httprouter"
)

func init() {
	module_manager.RegisterModule( new(DefaultWebServer) )
}

// Structures declaration
//-----------------------

const (
	CONFIG_PORT_KEY = "defaultWebServer.port"
)

type DefaultWebServer struct {
	Router *httprouter.Router
	Config config.I_Config
}



// I_Module interface implementation
//----------------------------------

func (it *DefaultWebServer) GetModuleName() string { return "core/defaultWebServer" }
func (it *DefaultWebServer) GetModuleDepends() []string { return make([]string, 0) }

func (it *DefaultWebServer) ModuleMakeSysInit() error { return nil }

func (it *DefaultWebServer) ModuleMakeConfig() error {
	if it.Config = config.GetConfig(); it.Config != nil {

		portValidator := func(val interface{}) (interface{}, bool) { newVal,ok := val.(int); return newVal,ok }
		if err :=  it.Config.RegisterItem( CONFIG_PORT_KEY, portValidator, 8000 ); err != nil { return err }

	} else {
		return errors.New("Can't get config instance")
	}

	return nil
}

func (it *DefaultWebServer) ModuleMakeInit() error {
	return web_server.RegisterWebServer(it)
}

func (it *DefaultWebServer) ModuleMakeVerify() error {
	it.Router = httprouter.New()

	return nil
}
func (it *DefaultWebServer) ModuleMakeLoad() error { return nil }
func (it *DefaultWebServer) ModuleMakeInstall() error { return nil }
func (it *DefaultWebServer) ModuleMakePostInstall() error { return nil }




// I_WebServer interface implementation
//-------------------------------------

func (it *DefaultWebServer) Run() error {
	listenOn := ":" + strconv.Itoa( it.Config.GetValue( CONFIG_PORT_KEY ).(int) )
	log.Println("starting Web Server on " + listenOn)
	http.ListenAndServe(listenOn, it.Router)

	return nil
}

func (it *DefaultWebServer) GetName() string {
	return "DefaultWebServer"
}

func (it *DefaultWebServer) RegisterHandler() string {
	return "DefaultWebServer"
}

func (it *DefaultWebServer) RegisterController(HTTPType string, Path string, CallbackFunc interface{}) error {
	if CallbackHandler, ok := CallbackFunc.(httprouter.Handle); ok {
		switch HTTPType {
		case "GET":
			it.Router.GET(Path, CallbackHandler)
		case "POST":
			it.Router.POST(Path, CallbackHandler)
		default:
			return errors.New("can't register HTTP protocol type " + HTTPType)
		}
	} else {
		return errors.New("callback handler should be of [github.com/julienschmidt/httprouter.Handle] type")
	}

	return nil
}
