package default_config

import (
	config "github.com/ottemo/foundation/config"
	db "github.com/ottemo/foundation/database"
)

func init() {
	instance := new(DefaultConfig)

	db.RegisterOnDatabaseStart( instance.Load )
	config.RegisterConfig( new(DefaultConfig) )
}
