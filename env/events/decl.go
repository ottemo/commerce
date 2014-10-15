package events

import (
	"github.com/ottemo/foundation/env"
)

type DefaultEventBus struct {
	listeners []env.F_EventListener
}
