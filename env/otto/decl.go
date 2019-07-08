package otto

import (
	"github.com/robertkrimen/otto"
	"github.com/ottemo/commerce/env"
	"sync"
)

const (
	ConstSessionKey = "script_id"

	ConstErrorModule = "env/otto"
	ConstErrorLevel = env.ConstErrorLevelService
)

var engine *ScriptEngine

// ScriptEngine is an implementer of InterfaceScriptEngine
type ScriptEngine struct {
	mutex sync.RWMutex
	mappings map[string]interface{}
	instances map[string]*Script
}

// ScriptEngine is an implementer of InterfaceScriptEngine
type Script struct {
	id string
	vm *otto.Otto
}
