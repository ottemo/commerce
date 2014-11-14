// Package "events" is a default implementation for "I_EventBus" interface.
package events

import (
	"github.com/ottemo/foundation/env"
)

// I_EventBus implementer class
type DefaultEventBus struct {
	listeners map[string][]env.F_EventListener
}
