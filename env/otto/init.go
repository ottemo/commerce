package otto

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/utils"
	"github.com/ottemo/commerce/env"
	"fmt"
)


func init() {
	engine := new(ScriptEngine)
	engine.mappings = make(map[string]interface{})

	engine.Set("json", utils.EncodeToJSONString)
	engine.Set("printf",  fmt.Sprintf)

	engine.Set("getModel", models.GetModel)
	engine.Set("getModels", models.GetDeclaredModels)

	env.RegisterScriptEngine("Otto", engine)
}
