package conf

import (
	db "github.com/ottemo/foundation/database"
	env "github.com/ottemo/foundation/env"
)

func init() {
	instance := new(DefaultConfig)

	db.RegisterOnDatabaseStart(instance.Load)
	env.RegisterConfig(new(DefaultConfig))
}
