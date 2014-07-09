package config

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

func init() {
	instance := new(DefaultConfig)

	db.RegisterOnDatabaseStart(instance.Load)
	env.RegisterConfig(new(DefaultConfig))
}
