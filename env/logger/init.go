package logger

import (
	"github.com/ottemo/foundation/env"
)

func init() {
	instance := new(DefaultLogger)
	var _ env.I_Logger = instance

	env.RegisterLogger(instance)
}
