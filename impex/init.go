package impex

import (
	"fmt"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	api.RegisterOnRestServiceStart(setupAPI)

	RegisterImportCommand("INSERT", new(ImportCmdInsert))
	RegisterImportCommand("UPDATE", new(ImportCmdUpdate))
	RegisterImportCommand("DELETE", new(ImportCmdDelete))

	RegisterImportCommand("STORE", new(ImportCmdStore))
	RegisterImportCommand("MEDIA", new(ImportCmdMedia))

	RegisterImportCommand("ATTRIBUTE_ADD", new(ImportCmdAttributeAdd))

	// initializing column conversion functions
	ConversionFuncs["log"] = func(args ...interface{}) {
		env.Log(ConstLogFileName, env.ConstLogPrefixDebug, fmt.Sprint(args))
	}

	ConversionFuncs["logf"] = func(format string, args ...interface{}) {
		env.Log(ConstLogFileName, env.ConstLogPrefixDebug, fmt.Sprintf(format, args))
	}

	ConversionFuncs["set"] = func(context map[string]interface{}, key string, value interface{}) string {
		context[key] = value
		return ""
	}
}
