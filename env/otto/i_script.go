package otto

import (
	"github.com/robertkrimen/otto/repl"
)

func (it *Script) GetId() string {
	return it.id
}

func (it *Script) Execute(code string) (interface{}, error) {
	value, err := it.vm.Run(code)
	if err != nil {
		return nil, err
	}
	return value.Export()
}

func (it *Script) Get(name string) (interface{}, error) {
	return it.vm.Get(name)
}

func (it *Script) Set(name string, value interface{}) error {
	return it.vm.Set(name, value)
}

func (it *Script) Interact() error {
	return repl.RunWithOptions(it.vm, repl.Options{ Prompt: "otto> ", Autocomplete: true })
}

