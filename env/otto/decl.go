package otto

import (
	"github.com/robertkrimen/otto"
	"sync"
	"github.com/ottemo/commerce/env"
)
const (
	ConstErrorModule = "env/otto"
	ConstErrorLevel = env.ConstErrorLevelService
)

// ScriptEngine is an implementer of InterfaceScriptEngine
type ScriptEngine struct {
	mutex sync.RWMutex
	mappings map[string]interface{}
	instances []*otto.Otto
}

// ScriptEngine is an implementer of InterfaceScriptEngine
type Script struct {
	vm *otto.Otto
}
