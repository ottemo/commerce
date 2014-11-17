// Package events is a default implementation of I_EventBus declared in
// "github.com/ottemo/foundation/env" package
package events

import (
	"github.com/ottemo/foundation/env"
)

// I_EventBus implementer class
type DefaultEventBus struct {
	listeners map[string][]env.F_EventListener
}
