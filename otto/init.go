package otto

import (
	"github.com/robertkrimen/otto"
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/utils"
	"fmt"
)


func init() {

	baseVM := otto.New()

	// baseVM.Set("model", otto.Object{})

	baseVM.Run(`
		function dir(subject, level) {
			prefix = '';
			for (var i=0; i<level; i++) {
				prefix += ' ';
			}
			
			for (var x in subject) {
				if (typeof x == 'object') {
					dump(x, level+1);
				} else {
					console.log(prefix + x);
				}
			}
                }
	`)

	baseVM.Set("json", utils.EncodeToJSONString)
	baseVM.Set("getModel", models.GetModel)
	baseVM.Set("models", models.GetDeclaredModels)
	baseVM.Set("print", func(x interface{}) string { return fmt.Sprintf("%o", x) })

	engine := new(ScriptEngine)
	engine.baseVM = baseVM

	models.RegisterScriptEngine("Otto", engine)
}
