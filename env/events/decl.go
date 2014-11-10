package events

import (
	"github.com/ottemo/foundation/env"
)

type DefaultEventBus struct {
	listeners map[string][]env.F_EventListener
}
