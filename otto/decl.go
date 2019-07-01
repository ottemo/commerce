package otto

import (
	"github.com/robertkrimen/otto"
)

// ScriptEngine is an implementer of InterfaceScriptEngine
type ScriptEngine struct {
	baseVM *otto.Otto
}

// ScriptEngine is an implementer of InterfaceScriptEngine
type Script struct {
	vm *otto.Otto
}
