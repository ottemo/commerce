package sqlite

import (
	database "github.com/ottemo/foundation/database"
	config "github.com/ottemo/foundation/config"
	sqlite3 "code.google.com/p/go-sqlite/go1/sqlite3"
)

func init() {
	instance := new(SQLite)

	config.RegisterOnConfigIniStart( instance.Startup )
	database.RegisterDBEngine( instance )
}


func (it *SQLite) Startup() error {

	var uri string = "ottemo.db"

	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.sqlite3.uri"); iniValue != "" {
			uri = iniValue
		}
	}

	if newConnection, err := sqlite3.Open(uri); err == nil {
		it.Connection = newConnection
	} else {
		return err
	}

	database.OnDatabaseStart()

	return nil
}
