package otto

import "github.com/ottemo/commerce/app/models"

func (it *ScriptEngine) GetScriptName() string {
	return "Otto"
}

func (it *ScriptEngine) GetScriptInstance() models.InterfaceScript {
	script := new(Script)
	script.vm = it.baseVM.Copy()
	return script
}
