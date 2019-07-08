package otto

import (
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
	"github.com/robertkrimen/otto"
)

func (it *ScriptEngine) GetScriptName() string {
	return "Otto"
}

func (it *ScriptEngine) GetScriptInstance(id string) env.InterfaceScript {

	if instance, present := it.instances[id]; present {
		return instance
	}

	script := new(Script)
	script.vm = otto.New()
	script.id = id

	it.mutex.Lock()
	defer it.mutex.Unlock()

	for key, value := range it.mappings {
		script.vm.Set(key, value)
	}

	// TODO: implement instances cleanup (lifetime based)
	it.instances[script.id] = script

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
