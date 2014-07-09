package sqlite

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

func init() {
	instance := new(SQLite)

	env.RegisterOnConfigIniStart(instance.Startup)
	db.RegisterDBEngine(instance)
}

func (it *SQLite) Startup() error {

	var uri string = "ottemo.db"

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.sqlite3.uri", uri); iniValue != "" {
			uri = iniValue
		}
	}

	if newConnection, err := sqlite3.Open(uri); err == nil {
		it.Connection = newConnection
	} else {
		return err
	}

	db.OnDatabaseStart()

	return nil
}
