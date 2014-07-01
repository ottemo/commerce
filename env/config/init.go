package config

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/database"
)

func init() {
	instance := new(DefaultConfig)

	database.RegisterOnDatabaseStart( instance.Load )
	env.RegisterConfig( new(DefaultConfig) )
}
