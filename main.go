package main

import (
	"github.com/ottemo/platform/interfaces/config"

	"github.com/ottemo/platform/tools/callback_chain"
	"github.com/ottemo/platform/tools/module_manager"
	"github.com/ottemo/platform/tools/web_server"

	_ "github.com/ottemo/platform/modules/db_sqlite"
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

	web_server.Run()
	config.GetConfig().Save()

	callback_chain.ExecuteCallbackChain("shutdown")
}
