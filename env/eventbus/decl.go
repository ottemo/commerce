package eventbus

import (
	"github.com/ottemo/commerce/env"
)

// Package global constants
const (
	ConstErrorModule = "env/eventbus"
	ConstErrorLevel  = env.ConstErrorLevelService
)

var eventBus *DefaultEventBus

// DefaultEvent InterfaceEvent implementer
type DefaultEvent struct {
	id string
	data map[string]interface{}
	propagate bool
}

// DefaultEventListener InterfaceEvent implementer
type DefaultEventListener struct {
	id string
	handler env.FuncEventHandler
	priority float64
}

// DefaultEventBus InterfaceEventBus implementer
type DefaultEventBus struct {
	initialized bool
	events map[string]string
	listeners map[string][]env.InterfaceEventListener
}
