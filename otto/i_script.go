package otto

import "github.com/robertkrimen/otto/repl"

func (it *Script) Interact() error {
	return repl.Run(it.vm)
}

func (it *Script) Execute(code string) (interface{}, error) {
	return it.vm.Eval(code)
}

func (it *Script) Get(name string) (interface{}, error) {
	return it.vm.Get(name)
}

func (it *Script) Set(name string, value interface{}) error {
	return it.vm.Set(name, value)
}
