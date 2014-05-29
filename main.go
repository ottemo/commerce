package main

import (
	// strict packages (hard dependency)
	"github.com/ottemo/platform/interfaces/config"
	"github.com/ottemo/platform/interfaces/web_server"

	"github.com/ottemo/platform/tools/callback_chain"
	"github.com/ottemo/platform/tools/module_manager"


	// optional packages (soft dependency)
	_ "github.com/ottemo/platform/modules/db_sqlite"
	_ "github.com/ottemo/platform/modules/web_server"
	_ "github.com/ottemo/platform/modules/config"
)

func init() {
	callback_chain.RegisterCallbackChain("startup",
		callback_chain.DefaultVoidFuncWithNoParamsValidator,
		callback_chain.DefaultChainVoidFuncWithNoParamsExecutor)

	callback_chain.RegisterCallbackChain("shutdown",
		callback_chain.DefaultVoidFuncWithNoParamsValidator,
		callback_chain.DefaultChainVoidFuncWithNoParamsExecutor)
}

func main() {
	callback_chain.ExecuteCallbackChain("startup")

	if err := module_manager.InitModules(); err != nil {
		panic("Can't init modules: " + err.Error() )
	}

	if webServer := web_server.GetWebServer(); webServer != nil {
		if err := webServer.Run(); err != nil {
			panic("Web Server Run() issues: " + err.Error())
		}
	} else {
		panic("No web server found")
	}

	config.GetConfig().Save()

	callback_chain.ExecuteCallbackChain("shutdown")
}
