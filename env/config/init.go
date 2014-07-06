package config

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/db"
)

func init() {
	instance := new(DefaultConfig)

	db.RegisterOnDatabaseStart( instance.Load )
	env.RegisterConfig( new(DefaultConfig) )
}
