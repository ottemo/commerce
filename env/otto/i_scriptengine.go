package otto

import (
	"github.com/ottemo/commerce/env"
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


func (it *ScriptEngine) Get(name string) (interface{}, error) {
	it.mutex.Lock()
	defer it.mutex.Unlock()

	if value, present := it.mappings[name]; present {
		return value, nil
	}

	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "df850b5d-bbf0-4e5b-a6df-1b4883320c09", "Key '" + name + "' does not exist")
}

func (it *ScriptEngine) Set(name string, value interface{}) error {
	it.mutex.Lock()
	defer it.mutex.Unlock()

	if value == nil {
		if _, present := it.mappings[name]; present {
			delete(it.mappings, name)
		}
	} else {
		it.mappings[name] = value
	}

	return nil
}
