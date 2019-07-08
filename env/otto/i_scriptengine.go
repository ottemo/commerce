package otto

import (
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
	"github.com/robertkrimen/otto"
)

func (it *ScriptEngine) GetScriptName() string {
	return "Otto"
}

func (it *ScriptEngine) GetScriptInstance() env.InterfaceScript {

	script := new(Script)
	script.vm = otto.New()

	it.mutex.Lock()
	defer it.mutex.Unlock()

	for key, value := range it.mappings {
		script.vm.Set(key, value)
	}

	it.instances = append(it.instances, script.vm)

	return script
}


func (it *ScriptEngine) Get(path string) (interface{}, error) {
	it.mutex.Lock()
	defer it.mutex.Unlock()

	return utils.MapGetPathValue(it.mappings, path)
}

func (it *ScriptEngine) Set(path string, value interface{}) error {
	it.mutex.Lock()
	defer it.mutex.Unlock()

	if value == nil {
		return utils.MapSetPathValue(it.mappings, path, value, true)
	} else {
		return utils.MapSetPathValue(it.mappings, path, value, false)
	}
}
