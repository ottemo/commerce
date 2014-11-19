// Package events is a default implementation of InterfaceEventBus declared in
// "github.com/ottemo/foundation/env" package
package events

import (
	"github.com/ottemo/foundation/env"
)

// DefaultEventBus InterfaceEventBus implementer class
type DefaultEventBus struct {
	listeners map[string][]env.FuncEventListener
}
